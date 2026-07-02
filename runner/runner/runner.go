// Package runner предоставляет реализацию паттерна "Runner" (воркер-пул)
// для асинхронного выполнения ограниченного или динамического набора задач
// с поддержкой обработки сигналов прерывания ОС.
package runner

import (
	"log"
	"os"
	"os/signal"
	"patterns/runner/task"
	"sync"
	"time"
)

// Runner представляет собой менеджер задач, который управляет
// пулом воркеров (горутин) и координирует их выполнение.
type Runner struct {
	// interrupt принимает системные сигналы прерывания (например, SIGINT).
	interrupt chan os.Signal
	// timeout используется для ограничения общего времени работы (в текущей логике не инициализирован).
	timeout <-chan time.Time
	// workersNum определяет фиксированное количество одновременно работающих воркеров.
	workersNum uint
	// tasks — буферизированный канал для очереди входящих задач.
	tasks chan task.Task
	// completed — канал для сбора результатов выполнения задач.
	completed chan task.Completed
}

// New создает и инициализирует новый экземпляр Runner.
// workersNum задает количество одновременно активных горутин.
// tasksNum определяет емкость (буфер) канала задач.
func New(workersNum uint, tasksNum int) *Runner {
	return &Runner{
		interrupt:  make(chan os.Signal),
		timeout:    make(<-chan time.Time),
		workersNum: workersNum,
		tasks:      make(chan task.Task, tasksNum),
		completed:  make(chan task.Completed),
	}
}

// AddTask добавляет одну или несколько задач в очередь выполнения.
// Метод заблокирует поток, если буфер канала tasks будет переполнен.
func (r *Runner) AddTask(tasks ...task.Task) {
	for _, task := range tasks {
		r.tasks <- task
	}
}

// addRoutine запускает одну фоновую горутину (воркера).
// Воркер берет одну задачу из канала tasks, выполняет ее и отправляет
// результат в канал completed. После этого горутина завершается.
func (r *Runner) addRoutine() {
	log.Println("Adding new routine")
	go func() {
		var task task.Task
		task, ok := <-r.tasks
		log.Println("===")
		if !ok {
			log.Println("Failed to run task, because task pipe is closed")
			return
		}
		// Выполнение задачи и отправка структуры результата
		r.completed <- task.Run()
		log.Println("Finishing routine")
	}()
}

// Run запускает главный цикл управления воркерами в отдельной горутине.
// Метод подписывается на системные прерывания и инициализирует стартовый пул воркеров.
// Возвращает nil и не блокирует основной поток выполнения.
// При получении сигнала прерывания вызывает wg.Done() для уведомления внешней группы ожидания.
func (r *Runner) Run(wg *sync.WaitGroup) error {
	// Подписка на сигнал прерывания (например, Ctrl+C)
	signal.Notify(r.interrupt, os.Interrupt)

	// Инициализация начального пула воркеров
	for i := 0; i < int(r.workersNum); i++ {
		r.addRoutine()
	}

	// Запуск главного управляющего цикла
	go func() {
	runnerLoop:
		for {
			select {
			case <-r.interrupt:
				log.Println("Stopping runner")
				wg.Done() // Уведомление об успешной остановке
				break runnerLoop
			case completed := <-r.completed:
				// Логирование результата выполнения задачи
				if completed.Err != nil {
					log.Printf("Task finished with error: %v", completed.Err.Error())
				} else {
					log.Println("Task finished with success")
				}
				// Паттерн динамического восполнения пула:
				// вместо завершения пула, взамен отработавшей горутины создается новая.
				r.addRoutine()
			}
		}
	}()

	return nil
}
