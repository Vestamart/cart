package errGroup

import (
	"context"
	"sync"
)

type ErrGroup struct {
	wg      sync.WaitGroup     // Отслеживание активных горутин
	errOnce sync.Once          // Для сохранения первой ошибки
	ctx     context.Context    // Контекст для отмены
	cancel  context.CancelFunc // Функция отмены контекста
	err     error              // Первая ошибка
}

func NewErrGroup(ctx context.Context) (*ErrGroup, context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	return &ErrGroup{ctx: ctx, cancel: cancel}, ctx
}

func (g *ErrGroup) Go(f func(ctx context.Context) error) {
	g.wg.Add(1)
	go func() {
		defer g.wg.Done()
		if err := f(g.ctx); err != nil {
			g.errOnce.Do(func() {
				g.err = err
				g.cancel()
			})
		}
	}()
}

func (g *ErrGroup) Wait() error {
	g.wg.Wait()
	g.cancel()
	return g.err
}
