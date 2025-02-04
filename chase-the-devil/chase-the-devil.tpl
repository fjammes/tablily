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

{{ .TablilyNotes }}
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
