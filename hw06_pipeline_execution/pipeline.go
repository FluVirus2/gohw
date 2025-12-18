package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func RunStage(in In, out Bi, done In) {
	defer close(out)

	for {
		select {
		case <-done:
			for range in {
			}
			return
		case v, ok := <-in:
			if !ok {
				return
			}
			select {
			case out <- v:
			case <-done:
				return
			}
		}
	}
}

func CancellableChan(in In, done In) Out {
	out := make(Bi)
	go RunStage(in, out, done)
	return out
}

func CancellableStage(stage Stage, done In) Stage {
	return func(in In) Out {
		return CancellableChan(stage(in), done)
	}
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	st := CancellableChan(in, done)

	for _, stage := range stages {
		st = CancellableStage(stage, done)(st)
	}

	return st
}
