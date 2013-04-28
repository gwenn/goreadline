goreadline
==========

Yet another Readline binding.

SetCompletionEntryFunction should be used to register an application-specific completion function.  
The default/filename completion is called when there is no application-specific match.

AddHistory ignores space and consecutive dups.  
ReadHistory ignores syscall.ENOENT error (meaning that the history file doesn't exist).  
AppendHistory creates the history file if it doesn't exist.  
GetHistory supports negative index to ease browsing the last history entries.

### Readline documentation:

http://cnswww.cns.cwru.edu/php/chet/readline/readline.html  
http://cnswww.cns.cwru.edu/php/chet/readline/history.html

http://www.thrysoee.dk/editline/  
http://www.cs.utah.edu/~bigler/code/libedit.html

### Similar projects:

https://github.com/bobappleyard/readline  
https://github.com/sbinet/go-readline  
https://github.com/shavac/readline  

http://code.google.com/p/go.crypto/ssh/terminal (semi official)  
\- https://github.com/sbinet/go-terminal (go.crypto/ssh/terminal forked)  
https://github.com/kless/terminal  
https://github.com/peterh/liner  
\- https://github.com/sbinet/liner  
https://github.com/edsrzf/fineline  
\- https://github.com/davecheney/fineline  
\- https://github.com/sbinet/fineline  
