# mkrq: make a radio question
# input:  on the first line, Qno (question number; >0) to begin numbering from;
#         then on each line: filename of choices list to include, then question text
# output: Go literal to put into the questions slice

BEGIN     { qno = -1 }
qno == -1 { qno = $1; next }

qno > 0 {
	fname = $1
	text = ""

	# for some reason the following fails in Gawk (which is BS and doesn't
	# adhere to the Awk book)
	"awk -f mkchoices.awk " fname "" | getline choices

	for(i = 2; i <= NF; i++)
		text = text " " $i

	sub("^ ", "", text)

	printf("\t&question{%d, \"%s\", true, nil, %s},\n", qno, text, choices)
	qno++
}
