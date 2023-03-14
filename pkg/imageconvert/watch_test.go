package imageconvert

/*
var DummyEntry = path.Entry{
	FileInfo:     nil,
	AbsolutePath: "/home/kmulvey/src/go/src/github.com/kmulvey/imageconvert/cmd/imageconvert/main.go",
	Children:     []path.Entry{},
}

func TestWaitTilFileWritesComplete(t *testing.T) {
	t.Parallel()

	var eventsIn = make(chan path.WatchEvent)
	var eventsOut = make(chan path.WatchEvent)

	go waitTilFileWritesComplete(eventsIn, eventsOut)

	var create = path.WatchEvent{Entry: DummyEntry, Op: 1}
	var write = path.WatchEvent{Entry: DummyEntry, Op: 2}

	for i := 0; i < 1000; i++ {
		if i == 0 {
			eventsIn <- create
		} else {
			eventsIn <- write
		}
	}

	go func() {
		for e := range eventsOut {
			assert.True(t, strings.HasSuffix(e.Entry.AbsolutePath, "watch.go"))
		}
	}()

	time.Sleep(time.Second * 2)
	close(eventsIn)
}
*/
