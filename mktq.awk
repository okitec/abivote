# mktq.awk; create a text question entry
# input:  on the first line, Qno (question number; >0) to begin numbering from;
#         then question text, one question per line
# output: Go literal to put into the questions slice

BEGIN     { qno = -1 }
qno == -1 { qno = $1; next }
qno > 0   { printf("\t&question{%d, \"%s\", false, nil, nil},\n", qno, $0); qno++ }

