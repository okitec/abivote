# mkradio.awk: make Go []Choice from list of options
# input: one option per line
# output: one long line for the []Choice

# note the lack of newlines
BEGIN { printf("[]Choice{") }
      { printf("Choice{\"%s\", nil, 0.0}, ", $0) }
END   { printf("}\n") }
