package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
    "strings"
    // "time"
    "flag"
    "strconv"
    "regexp"
)

func main() {

  inPath := flag.String("in", ".", "path to input file")
  messageFieldNumbers := flag.String("fields", "0", "column positions in the message to output (differs between message types).")
  urnPosition := flag.Int("urnpos", 2, "column position of URN in file")
  outPath := flag.String("out", "./out.txt", "output file path")
  overwrite := flag.Bool("overwrite", false, "should the output be overwritten? If false, the output will append")
  delim := flag.String("delimiter", "\t", "output file delimiter")
  prefixData := flag.String("prefix", "", "Data to prefix to every line")
  linestoParse := flag.Int("lines", -1, "Number of lines to parse. default of -1 is all lines in file")

  flag.Parse()

  fmt.Println("Input path: ", *inPath)
  fmt.Println("Output path: ", *outPath)
  fmt.Println("Message fields: ", *messageFieldNumbers)

  // Open the input file
  file, err := os.Open(*inPath)
  if err != nil {
    log.Fatal(err)
  }
  defer file.Close()

  // Slice out the positions
  positions, _ := sliceAtoi(strings.Split(*messageFieldNumbers, ","))


  // Set up a slice to store the string output
  var strs []string
  scanner := bufio.NewScanner(file)

  cnt := 0
  for scanner.Scan() {
    // Split the line into '+'-separated chunks
    s := scanner.Text()
    var re = regexp.MustCompile(`^(\d{4}-\d{2}-\d{2}\s[\d:\.]*)\s`)
    s = re.ReplaceAllString(s, `$1+`)
    pieces := strings.Split(s, "+")

    // Work through the messageField Numbers and pull the elements from the
    //  input that correspond

    if len(pieces) > MaxIntSlice(positions) {
      var str []string
      for _, ind := range positions {
        str = append(str, pieces[ind]) // strip tabs
      }
      // fmt.Println(" --- Extracted: ", str)
      // Per the instruction, the first item in the output string should be the
      //  SupporterURN positions. If this is empty, add the str to the output
      if len(pieces[*urnPosition]) == 0 {
        strs = append(strs, fmt.Sprintf("%s%s%s", *prefixData, *delim, strings.Join(str, *delim)))
      } else {
        fmt.Println(" --- Ignoring: Isp ", str, " with URN ", pieces[*urnPosition])
      }
    }
    cnt++
    if *linestoParse > -1 && cnt >= *linestoParse {
      break;
    }
  }
  defer file.Close()

  fmt.Println("Lines to add: ", len(strs))

  // Create the output OpenFile
  createFile(*outPath, *overwrite)
  writeFile(*outPath, strs)

  if err := scanner.Err(); err != nil {
    log.Fatal(err)
  }
}

// Turns a slice ostrings into a slice of ints
func sliceAtoi(sa []string) ([]int, error) {
    si := make([]int, 0, len(sa))
    for _, a := range sa {
        i, err := strconv.Atoi(a)
        if err != nil {
            return si, err
        }
        si = append(si, i)
    }
    return si, nil
}

// Finds the largest value in a slice of ints
func MaxIntSlice(v []int) (m int) {
    if len(v) > 0 {
        m = v[0]
    }
    for i := 1; i < len(v); i++ {
        if v[i] > m {
            m = v[i]
        }
    }
    return
}


func createFile(path string, overwrite bool) {
	// detect if file exists
	var _, err = os.Stat(path)

	// create file if not exists
	if overwrite || os.IsNotExist(err) {
		var file, err = os.Create(path)
		if isError(err) { return }
		defer file.Close()
	}

	fmt.Println("==> done creating file", path)
}

func writeFile(path string, vals []string) {
	// open file using READ & WRITE permission
  var fileFlags = os.O_WRONLY|os.O_APPEND
	var file, err = os.OpenFile(path, fileFlags, 0644)
	if isError(err) { return }
	defer file.Close()

	// write some text line-by-line to file
  for _, str := range vals {
	  _, err = file.WriteString(str + "\r\n")
	  if isError(err) { return }
  }

	// save changes
	err = file.Sync()
	if isError(err) { return }

	fmt.Println("==> done writing to file")
}

func isError(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
	}
	return (err != nil)
}
