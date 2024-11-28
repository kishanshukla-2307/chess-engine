package chessBots

type Bot interface {
	BestMove()
}

type RandomBot struct {
}

func (rb *RandomBot) BestMove() {

}
