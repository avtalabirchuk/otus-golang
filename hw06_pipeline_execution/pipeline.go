package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

// in -> st1 -> done_stage -> st2 -> done_stage -> st3 -> done_stage -> out

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := doneStage(in, done)
	for _, stage := range stages {
		out = doneStage(stage(out), done)
	}

	return doneStage(out, done)
}

func doneStage(in In, done In) Out {
	out := make(Bi)

	go func() {
		defer close(out)
		for {
			select {
			case d, ok := <-in:
				if !ok {
					return
				}
				select {
				case out <- d:
				case <-done:
					return
				}
			case <-done:
				return
			}
		}
	}()
	return out
}
