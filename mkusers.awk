# mkusers.awk: produce users.json from list of user IDs

BEGIN  { printf("{") }
NF > 0 { printf("\"%s\":{\"Name\":\"%s\",\"HasVoted\":false,\"Admin\":false},\n", $0, $0) }
END    { printf("}\n") }
