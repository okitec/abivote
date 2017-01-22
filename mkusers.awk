# mkusers.awk: produce users.json from list of user IDs
# ATTENZIONE: Remove the trailing comma in the last line because JSON.
# lest you trigger 'invalid character '}' looking for beginning of object key string'

BEGIN  { printf("{") }
NF > 0 { printf("\"%s\":{\"Name\":\"%s\",\"HasVoted\":false,\"Admin\":false},\n", $0, $0) }
END    { printf("}\n") }
