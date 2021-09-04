package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

// in -> st1 -> doneStage -> st2 -> doneStage -> st3 -> doneStage -> out

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	if in == nil {
		return nil
	}

	out := in
	for _, stage := range stages {
		out = stage(doneStage(out, done))
	}

	return doneStage(out, done)
}

func doneStage(in In, done In) Out {
	out := make(Bi)

	go func() {
		defer func() {
			close(out)
			for range in {
			}
		}()
		for {
			select {
			case d, ok := <-in:
				if !ok {
					return
				}
				out <- d
			case <-done:
				return
			}
		}
	}()
	return out
}
