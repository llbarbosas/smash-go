package state

type State struct {
	Board      []uint8 `json:"board"`
	NextPlayer uint8   `json:"next_player"`
}

func (s *State) Update(move uint8) {
	if s.Board[move] == 0 {
		s.Board[move] = s.NextPlayer
	}

	if s.NextPlayer == 1 {
		s.NextPlayer = 2
	} else {
		s.NextPlayer = 1
	}
}

func (s *State) HaveWinner() uint8 {
	diagonalWin := s.Board[4] != 0 && ((s.Board[0] == s.Board[4] && s.Board[4] == s.Board[8]) || (s.Board[2] == s.Board[4] && s.Board[4] == s.Board[6]))

	if diagonalWin {
		return s.Board[4]
	}

	for i := 0; i < 3; i++ {
		if s.Board[i*3] != 0 {
			rowWin := s.Board[i*3] == s.Board[i*3+1] && s.Board[i*3+1] == s.Board[i*3+2]
			colWin := s.Board[i*3] == s.Board[(i+1)*3] && s.Board[(i+1)*3] == s.Board[(i+2)*3]

			if rowWin || colWin {
				return s.Board[i*3]
			}
		}
	}

	return 0
}
