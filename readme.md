#H1 Tigre is a basic KeyDB / golang interloper

**More info on installing KeyDB here -- https://github.com/Snapchat/KeyDB**

#H2 - KDB 
- Connect to a KeyDB (*redis.Client)
- List functionality:  
    Txt2List        -> convert a txt file into a KDB list
    List_RandItem   -> Get random key from list
    List_Add        -> Add key to list / LPUSH
    List_DelStr     -> Delete key from list / LREM
    List2JSON       -> Return list as JSON / LRANGE 0, -1 
    List2JSON_alpha -> Return list as JSON & []string

#H2 - API client (main.go)

-- toy examples --
"/jFruit" -  KDB.List2JSON(keyDB, fruit_key)
"/jGods"  -  KDB.List2JSON(keyDB, gods_key)

