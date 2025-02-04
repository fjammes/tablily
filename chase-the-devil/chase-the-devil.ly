\version "2.24.3"

\header {
  title = "Chase the Devil"
  composer = "Max Romeo"
}


bassTab = \relative c, {
  \clef "bass_8"
  \key c \major
  \time 4/4
  \tempo 4 = 92

e4\3 e,4\4 r4 r4
r1
e'4\3 e,4\4 g4\4

}


\score {
  <<
    \new Staff {
      \bassTab
    }
    \new TabStaff \with {
    stringTunings = #bass-tuning
  } {
      \bassTab
    }
  >>
}
