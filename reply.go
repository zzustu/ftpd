package ftpd

const (
	// 110 Restart marker reply. In this case, the text is exact and not left to
	// the particular implementation; it must read: MARK yyyy = mmmm Where yyyy
	// is User-process data stream marker, and mmmm server's equivalent marker
	// (note the spaces between markers and "=").
	reply110RestartMarkerreply = 110

	// 120 Service ready in nnn minutes.
	reply120ServiceReadyInNnnMinutes = 120

	// 125 Data connection already open; transfer starting.
	reply125DataConnectionAlreadyOpen = 125

	// 150 File status okay; about to open data connection.
	reply150FileStatusOkay = 150

	// 200 Command okay.
	reply200CommandOkay = 200

	// 202 Command not implemented, superfluous at this site.
	reply202CommandNotImplemented = 202

	// 211 System status, or system help reply.
	reply211SystemStatusreply = 211

	// 212 Directory status.
	reply212DirectoryStatus = 212

	// 213 File status.
	reply213FileStatus = 213

	// 214 Help message. On how to use the server or the meaning of a particular
	// non-standard command. This reply is useful only to the human user.
	reply214HelpMessage = 214

	// 215 Name system type. Where Name is an official system name from the list
	// in the Assigned Numbers document.
	reply215NameSystemType = 215

	// 220 Service ready for new user.
	reply220ServiceReady = 220

	// Service closing control connection. Logged out if appropriate.
	reply221ClosingControlConnection = 221

	// 225 Data connection open; no transfer in progress.
	reply225DataConnectionOpenNoTransferInProgress = 225

	// Closing data connection. Requested file action successful (for example,
	// file transfer or file abort).
	reply226ClosingDataConnection = 226

	// 227 Entering Passive Mode (h1,h2,h3,h4,p1,p2).
	reply227EnteringPassiveMode = 227

	// 230 User logged in, proceed.
	reply230UserLoggedIn = 230

	// 250 Requested file action okay, completed.
	reply250RequestedFileActionOkay = 250

	// 257 "PATHName" created.
	reply257PathNameCreated = 257

	// 331 User name okay, need password.
	reply331UserNameOkayNeedPassword = 331

	// 332 Need account for login.
	reply332NeedAccountForLogin = 332

	// 350 Requested file action pending further information.
	reply350RequestedFileActionPendingFurtherInformation = 350

	// 421 Service not available, closing control connection. This may be a
	// reply to any command if the service knows it must shut down.
	reply421ServiceNotAvailableClosingControlConnection = 421

	// 425 Can't open data connection.
	reply425CantOpenDataConnection = 425

	// 426 Connection closed; transfer aborted.
	reply426ConnectionClosedTransferAborted = 426

	// 450 Requested file action not taken. File unavailable (e.g., file busy).
	reply450RequestedFileActionNotTaken = 450

	// 451 Requested action aborted: local error in processing.
	reply451RequestedActionAborted = 451

	// 452 Requested action not taken. Insufficient storage space in system.
	reply452RequestedActionNotTaken = 452

	// 500 Syntax error, command unrecognized. This may include errors such as
	// command line too long.
	reply500SyntaxErrorCommandUnrecognized = 500

	// 501 Syntax error in parameters or arguments.
	reply501SyntaxErrorInParametersOrArguments = 501

	// 502 Command not implemented.
	reply502CommandNotImplemented = 502

	// 503 Bad sequence of commands.
	reply503BadSequenceOfCommands = 503

	// 504 Command not implemented for that parameter.
	reply504CommandNotImplementedForThatParameter = 504

	// 530 Not logged in.
	reply530NotLoggedIn = 530

	// 532 Need account for storing files.
	reply532NeedAccountForStoringFiles = 532

	// 550 Requested action not taken. File unavailable (e.g., file not found,
	// no access).
	reply550RequestedActionNotTaken = 550

	// 551 Requested action aborted: page type unknown.
	reply551RequestedActionAbortedPageTypeUnknown = 551

	// 552 Requested file action aborted. Exceeded storage allocation (for
	// current directory or dataset).
	reply552RequestedFileActionAbortedExceededStorage = 552

	// 553 Requested action not taken. File name not allowed.
	reply553RequestedActionNotTakenFileNameNotAllowed = 553
)
