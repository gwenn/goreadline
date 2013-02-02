goreadline
==========

Yet another Readline binding.

Readline returns a (string, os.Error) couple like Reader#Read.

SetCompletionEntryFunction should be used to register an application-specific completion function.  
The default/filename completion is called when there is no application-specific match.

AddHistory ignores space and consecutive dups.  
ReadHistory ignores syscall.ENOENT error (meaning that the history file doesn't exist).  
AppendHistory creates the history file if it doesn't exist.  
GetHistory supports negative index to ease browsing the last history entries.  