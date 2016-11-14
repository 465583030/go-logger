package logger

import (
	"fmt"
	"testing"

	"os"

	"path/filepath"

	assert "github.com/blendlabs/go-assert"
)

func TestNewFileWriterUncompressed(t *testing.T) {
	assert := assert.New(t)

	tempFile := UUIDv4()
	uncompressed, err := NewFileWriter(tempFile, false, FileWriterUnlimitedSize, FileWriterUnlimitedArchiveFiles)
	assert.Nil(err)
	assert.NotNil(uncompressed)
	defer uncompressed.Close()
	defer func() {
		os.Remove(tempFile)
	}()

	assert.False(uncompressed.shouldCompressArchivedFiles)
	assert.Equal(0, uncompressed.fileMaxSizeBytes)
	assert.Equal(0, uncompressed.fileMaxArchiveCount)
	assert.Equal(fmt.Sprintf(isArchiveFileRegexpFormat, filepath.Base(tempFile)), uncompressed.isArchiveFileRegexp.String())

	stat, err := uncompressed.file.Stat()
	assert.Nil(err)
	assert.NotNil(stat)
	assert.Equal(filepath.Base(tempFile), stat.Name())

	written, err := uncompressed.Write([]byte(`this is only a test`))
	assert.Nil(err)
	assert.Equal(19, written)
}

func TestNewFileWriterCompressed(t *testing.T) {
	assert := assert.New(t)

	tempFile := UUIDv4()
	compressed, err := NewFileWriter(tempFile, true, FileWriterUnlimitedSize, FileWriterUnlimitedArchiveFiles)
	assert.Nil(err)
	assert.NotNil(compressed)
	defer compressed.Close()
	defer func() {
		os.Remove(tempFile)
	}()

	assert.True(compressed.shouldCompressArchivedFiles)
	assert.Equal(0, compressed.fileMaxSizeBytes)
	assert.Equal(0, compressed.fileMaxArchiveCount)
	assert.Equal(fmt.Sprintf(isCompressedArchiveFileRegexpFormat, filepath.Base(tempFile)), compressed.isArchiveFileRegexp.String())

	stat, err := compressed.file.Stat()
	assert.Nil(err)
	assert.NotNil(stat)
	assert.Equal(filepath.Base(tempFile), stat.Name())

	written, err := compressed.Write([]byte(`this is only a test`))
	assert.Nil(err)
	assert.Equal(19, written)
}

func TestNewFileWriterCompressedArchived(t *testing.T) {
	assert := assert.New(t)

	tempFile := UUIDv4()
	archived, err := NewFileWriter(tempFile, true, Kilobyte, FileWriterUnlimitedArchiveFiles)
	assert.Nil(err)
	assert.NotNil(archived)
	defer archived.Close()
	defer func() {
		os.Remove(tempFile)
	}()

	assert.True(archived.shouldCompressArchivedFiles)
	assert.Equal(Kilobyte, archived.fileMaxSizeBytes)
	assert.Equal(0, archived.fileMaxArchiveCount)
	assert.Equal(fmt.Sprintf(isCompressedArchiveFileRegexpFormat, filepath.Base(tempFile)), archived.isArchiveFileRegexp.String())

	stat, err := archived.file.Stat()
	assert.Nil(err)
	assert.NotNil(stat)
	assert.Equal(filepath.Base(tempFile), stat.Name())

	var written, total int
	for total < int(4*Kilobyte) {
		written, err = archived.Write([]byte("this is only a test\n"))
		assert.Nil(err)
		total += written
	}

	files, err := archived.getArchivedFilePaths()
	assert.Nil(err)
	defer func() {
		for _, path := range files {
			os.Remove(path)
		}
	}()
	assert.Len(files, 3)
}

func TestFileWriterShiftArchivedFiles(t *testing.T) {
	assert := assert.New(t)
	var err error

	td, err := os.Getwd()
	assert.Nil(err)
	id := UUIDv4()

	f1 := filepath.Join(td, id+".1")
	f2 := filepath.Join(td, id+".2")
	f3 := filepath.Join(td, id+".3")
	f4 := filepath.Join(td, id+".4")
	f5 := filepath.Join(td, id+".5")

	err = File.CreateAndClose(f1)
	assert.Nil(err)
	err = File.CreateAndClose(f2)
	assert.Nil(err)
	err = File.CreateAndClose(f3)
	assert.Nil(err)
	err = File.CreateAndClose(f4)
	assert.Nil(err)

	defer File.RemoveMany(f2, f3, f4, f5)

	regex, err := createIsArchivedFileRegexp(filepath.Join(td, id))
	assert.Nil(err)

	fr := &FileWriter{
		filePath:            filepath.Join(td, id),
		isArchiveFileRegexp: regex,
	}

	err = fr.shiftArchivedFiles([]string{f1, f2, f3, f4})
	assert.Nil(err)

	results, err := File.List(td, regex)
	assert.Nil(err)
	assert.Len(results, 4)
	assert.Equal(filepath.Join(td, id+".2"), results[0])
	assert.Equal(filepath.Join(td, id+".3"), results[1])
	assert.Equal(filepath.Join(td, id+".4"), results[2])
	assert.Equal(filepath.Join(td, id+".5"), results[3])
}

func TestFileWriterShiftCompressedArchivedFiles(t *testing.T) {
	assert := assert.New(t)
	var err error

	td, err := os.Getwd()
	assert.Nil(err)
	id := UUIDv4()

	f1 := filepath.Join(td, id+".1.gz")
	f2 := filepath.Join(td, id+".2.gz")
	f3 := filepath.Join(td, id+".3.gz")
	f4 := filepath.Join(td, id+".4.gz")
	f5 := filepath.Join(td, id+".5.gz")

	err = File.CreateAndClose(f1)
	assert.Nil(err)
	err = File.CreateAndClose(f2)
	assert.Nil(err)
	err = File.CreateAndClose(f3)
	assert.Nil(err)
	err = File.CreateAndClose(f4)
	assert.Nil(err)

	defer File.RemoveMany(f2, f3, f4, f5)

	regex, err := createIsCompressedArchiveFileRegexp(filepath.Join(td, id))
	assert.Nil(err)

	fr := &FileWriter{
		filePath:                    filepath.Join(td, id),
		shouldCompressArchivedFiles: true,
		isArchiveFileRegexp:         regex,
	}

	err = fr.shiftArchivedFiles([]string{f1, f2, f3, f4})
	assert.Nil(err)

	results, err := File.List(td, regex)
	assert.Nil(err)
	assert.Len(results, 4)
	assert.Equal(filepath.Join(td, id+".2.gz"), results[0])
	assert.Equal(filepath.Join(td, id+".3.gz"), results[1])
	assert.Equal(filepath.Join(td, id+".4.gz"), results[2])
	assert.Equal(filepath.Join(td, id+".5.gz"), results[3])
}

func TestFileWriterShiftCompressedArchivedFilesWithMax(t *testing.T) {
	assert := assert.New(t)
	var err error

	td, err := os.Getwd()
	assert.Nil(err)
	id := UUIDv4()

	f1 := filepath.Join(td, id+".1.gz")
	f2 := filepath.Join(td, id+".2.gz")
	f3 := filepath.Join(td, id+".3.gz")
	f4 := filepath.Join(td, id+".4.gz")

	err = File.CreateAndClose(f1)
	assert.Nil(err)
	err = File.CreateAndClose(f2)
	assert.Nil(err)
	err = File.CreateAndClose(f3)
	assert.Nil(err)
	err = File.CreateAndClose(f4)
	assert.Nil(err)

	defer File.RemoveMany(f2, f3)

	regex, err := createIsCompressedArchiveFileRegexp(filepath.Join(td, id))
	assert.Nil(err)

	fr := &FileWriter{
		filePath:                    filepath.Join(td, id),
		shouldCompressArchivedFiles: true,
		isArchiveFileRegexp:         regex,
		fileMaxArchiveCount:         3,
	}

	err = fr.shiftArchivedFiles([]string{f1, f2, f3, f4})
	assert.Nil(err)

	results, err := File.List(td, regex)
	assert.Nil(err)
	assert.Len(results, 2, fmt.Sprintf("%#v", results))
	assert.Equal(filepath.Join(td, id+".2.gz"), results[0])
	assert.Equal(filepath.Join(td, id+".3.gz"), results[1])
}

func TestFileWriterExtractArchivedFileIndex(t *testing.T) {
	assert := assert.New(t)

	regex, err := createIsArchivedFileRegexp("stdout")
	assert.Nil(err)

	fw := &FileWriter{
		isArchiveFileRegexp: regex,
	}

	index, err := fw.extractArchivedFileIndex("stdout.1")
	assert.Nil(err)
	assert.Equal(1, index)

	index, err = fw.extractArchivedFileIndex("stdout.22")
	assert.Nil(err)
	assert.Equal(22, index)

	index, err = fw.extractArchivedFileIndex("stdout")
	assert.NotNil(err)
	assert.Equal(0, index)
}
