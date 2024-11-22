package main

import (
	"fmt"

	myChessEngine "github.com/kishanshukla-2307/chess-engine"
)

func main() {
	engine, err := myChessEngine.NewNoobEngine(false)
	if err != nil {
		fmt.Println(fmt.Errorf(err.Error()))
	}
	// 	engine.PrintBoard()
	// 	err = engine.MakeMove(myChessEngine.WHITE_PAWN, &myChessEngine.Pos{1, 4}, &myChessEngine.Pos{3, 4})
	// 	if err != nil {
	// 		fmt.Println(fmt.Errorf(err.Error()))
	// 	}
	// 	err = engine.MakeMove(myChessEngine.BLACK_PAWN, &myChessEngine.Pos{6, 4}, &myChessEngine.Pos{4, 4})
	// 	if err != nil {
	// 		fmt.Println(fmt.Errorf(err.Error()))
	// 	}
	// 	err = engine.MakeMove(myChessEngine.WHITE_PAWN, &myChessEngine.Pos{1, 3}, &myChessEngine.Pos{3, 3})
	// 	if err != nil {
	// 		fmt.Println(fmt.Errorf(err.Error()))
	// 	}
	// 	err = engine.MakeMove(myChessEngine.BLACK_PAWN, &myChessEngine.Pos{4, 4}, &myChessEngine.Pos{3, 3})
	// 	if err != nil {
	// 		fmt.Println(fmt.Errorf(err.Error()))
	// 	}
	// 	err = engine.MakeMove(myChessEngine.WHITE_KNIGHT, &myChessEngine.Pos{0, 6}, &myChessEngine.Pos{2, 5})
	// 	if err != nil {
	// 		fmt.Println(fmt.Errorf(err.Error()))
	// 	}
	// 	err = engine.MakeMove(myChessEngine.BLACK_KNIGHT, &myChessEngine.Pos{7, 6}, &myChessEngine.Pos{5, 5})
	// 	if err != nil {
	// 		fmt.Println(fmt.Errorf(err.Error()))
	// 	}
	// 	err = engine.MakeMove(myChessEngine.WHITE_PAWN, &myChessEngine.Pos{1, 0}, &myChessEngine.Pos{3, 0})
	// 	if err != nil {
	// 		fmt.Println(fmt.Errorf(err.Error()))
	// 	}
	// 	err = engine.MakeMove(myChessEngine.BLACK_PAWN, &myChessEngine.Pos{6, 7}, &myChessEngine.Pos{4, 7})
	// 	if err != nil {
	// 		fmt.Println(fmt.Errorf(err.Error()))
	// 	}
	// 	err = engine.MakeMove(myChessEngine.WHITE_BISHOP, &myChessEngine.Pos{0, 5}, &myChessEngine.Pos{3, 2})
	// 	if err != nil {
	// 		fmt.Println(fmt.Errorf(err.Error()))
	// 	}
	// 	err = engine.MakeMove(myChessEngine.BLACK_BISHOP, &myChessEngine.Pos{7, 5}, &myChessEngine.Pos{4, 2})
	// 	if err != nil {
	// 		fmt.Println(fmt.Errorf(err.Error()))
	// 	}
	// 	err = engine.MakeMove(myChessEngine.WHITE_BISHOP, &myChessEngine.Pos{0, 2}, &myChessEngine.Pos{4, 6})
	// 	if err != nil {
	// 		fmt.Println(fmt.Errorf(err.Error()))
	// 	}
	// 	err = engine.MakeMove(myChessEngine.BLACK_KNIGHT, &myChessEngine.Pos{7, 1}, &myChessEngine.Pos{5, 2})
	// 	if err != nil {
	// 		fmt.Println(fmt.Errorf(err.Error()))
	// 	}
	// 	err = engine.MakeMove(myChessEngine.WHITE_BISHOP, &myChessEngine.Pos{3, 2}, &myChessEngine.Pos{6, 5})
	// 	if err != nil {
	// 		fmt.Println(fmt.Errorf(err.Error()))
	// 	}
	// 	err = engine.MakeMove(myChessEngine.BLACK_PAWN, &myChessEngine.Pos{3, 3}, &myChessEngine.Pos{2, 3})
	// 	if err != nil {
	// 		fmt.Println(fmt.Errorf(err.Error()))
	// 	}
	// 	err = engine.MakeMove(myChessEngine.WHITE_KING, &myChessEngine.Pos{0, 4}, &myChessEngine.Pos{1, 3})
	// 	if err != nil {
	// 		fmt.Println(fmt.Errorf(err.Error()))
	// 	}

	if err := engine.Run(); err != nil {
		fmt.Println(err)
	}
}
