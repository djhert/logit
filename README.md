# Logit

Logger created in Golang.  Runs as a separate goroutine.

# Usage
Create a Loggit Object in your code, this will need to be a global variable so all functions can access it

	var ThisLog *Logit.Logger

Logit accepts items that implement the "io.WriteCloser" to be valid Loggers.  Included are funtions to Log to a File and the Terminal.

	FileLog, err := Logit.OpenFile("/path/to/log/file")
	TermLog := Logit.TermLog()


Somewhere in your code, you need to Start the Logging goroutine.

	ThisLog, err = Logit.StartLogger(FileLog)
	ThatLog, err = Logit.StartLogger(TermLog)

This will start a goroutine that has the file open.

Logit then uses Channels to push strings into this goroutine, writing them to the Log.
Format a string to your liking, and simply pass it into the Log object that you created.

	ThisLog.Log("string!")

Using the above will produce the following example output in the log:

	09:15:15.918788110:  string!

On shutdown of your application, you will want to run the "Quit" function on your Logit Object to ensure all is closed and the goroutines are stopped.

	ThisLog.Quit()
	ThatLog.Quit()

# Features
Able to handle concurrency.  As long as your goroutine has access to a valid Logit Object, you are able to use the Log function safely.
By pushing a string onto the channel, Logit will append it to the log when it can.

# TODO
As of now if the log file is removed during runtime Logit will not complain, and all will continue to work without a hitch except logging.  This is not ideal of course.
