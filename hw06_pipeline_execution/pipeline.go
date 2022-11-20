package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := in
	for _, stage := range stages {
		out = trackStage(done, stage(in))
		in = out
	}
	return out
}

func trackStage(done, in In) Out {
	toNextStage := make(Bi)

	go func() {
		defer finalizeStage(toNextStage, in)
		for {
			select {
			case <-done:
				return
			case v, ok := <-in:
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
