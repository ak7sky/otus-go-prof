package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	var out Out
	for _, stage := range stages {
		out = runStage(done, in, stage)
		in = out
	}
	return out
}

func runStage(done, in In, stage Stage) Out {
	toNextStage := make(Bi)
	out := stage(in)

	go func() {
		defer finalizeStage(toNextStage, in)
		for {
			select {
			case <-done:
				return
			case v, ok := <-out:
				if !ok {
					return
				}
				toNextStage <- v
			}
		}
	}()

	return toNextStage
}

func finalizeStage(toNextStage Bi, in In) {
	close(toNextStage)
	go func() {
		for range in {
		}
	}()
}
