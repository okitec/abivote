# mkq: make a question
# input: each line starts with "text" for text questions or "radio" for radio
#        questions; the second field on radio questions specifies the option
#        file to include.
# example:
#     text Who are you?
#     radio yes-no Do jellyfish dream?
# output: Go literal to put into the questions slice

BEGIN {
	qno = 1
}

$1 == "radio" {
	fname = $2

	# for some reason the following fails in Gawk (which is BS and doesn't
	# adhere to the Awk book)
	"awk -f mkchoices.awk " fname "" | getline choices

	text = catflds(3, NF)
	printf("\t&question{%d, \"%s\", true, nil, %s, \"\"},\n", qno, text, choices)
	qno++
}

$1 == "text" {
	text = catflds(2, NF)
	printf("\t&question{%d, \"%s\", false, nil, nil, \"\"},\n", qno, text)
	qno++
}

# catflds: concatenate fields $s to $e
function catflds(s, e,            text, i) {
	text = ""

	for(i = s; i <= e; i++)
		text = text " " $i

	sub("^ ", "", text)
	return text
}
