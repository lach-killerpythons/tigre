#H1 Tigre is a basic KeyDB endpoint API - **GET POST DELETE**

**More info on installing KeyDB here -- https://github.com/Snapchat/KeyDB**

#H2 - KDB 
- Connect to a KeyDB (*redis.Client) // *DB
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

-- prototype endpoints --
"/anylist" - jsonWildtype2 
1. Parse the "list" value from URL to listName  
` /anylist?list=fruits ` -> fruits

2. KDB.List2JSON_alpha(*DB, listName) => listObject (listname, []list)

3. write byte array to JSON & save listObj containing []string of list

"/new_god" - postGod

"/del" - delGod