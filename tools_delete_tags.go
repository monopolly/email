package email

func deleteTags(html []byte, tagStart, tagEnd string) []byte {

	lenStart := len(tagStart)
	lenEnd := len(tagEnd)

	var part, end []byte
	var begin bool
	var start, ends int

	for pos, x := range html {
		//fmt.Println(pos, string(x))
		switch begin {
		case false:
			//not found
			part = append(part, x)
			if len(part) == lenStart+1 {
				part = part[1:]
			}
			if string(part) == tagStart {
				//		fmt.Println("start", pos-lenStart+1, string(html[pos-lenStart+1]))
				//1 2 3 4 | 5
				start = pos - lenStart + 1
				begin = true
			}
			continue
		case true:
			//found
			end = append(end, x)
			if len(end) == lenEnd+1 {
				end = end[1:]
			}
			if string(end) == tagEnd {
				//		fmt.Println("end", pos, string(html[pos]))
				ends = pos
				begin = false

				//		fmt.Println("need delete", start, ends)
				for x := start; x < ends+1; x++ {
					html[x] = ' '
				}
				continue
			}
		}

	}
	return html
}

func deleteTag(html []byte, tagStart, tagEnd string) []byte {

	lenStart := len(tagStart)
	lenEnd := len(tagEnd)

	var part, end []byte
	var begin bool
	var start, ends int

	for pos, x := range html {
		//fmt.Println(pos, string(x))
		switch begin {
		case false:
			//not found
			part = append(part, x)
			if len(part) == lenStart+1 {
				part = part[1:]
			}
			if string(part) == tagStart {
				//		fmt.Println("start", pos-lenStart+1, string(html[pos-lenStart+1]))
				//1 2 3 4 | 5
				start = pos - lenStart + 1
				begin = true
			}
			continue
		case true:
			//found
			end = append(end, x)
			if len(end) == lenEnd+1 {
				end = end[1:]
			}
			if string(end) == tagEnd {
				//		fmt.Println("end", pos, string(html[pos]))
				ends = pos
				begin = false

				//		fmt.Println("need delete", start, ends)
				for x := start; x < ends+1; x++ {
					html[x] = ' '
				}
				return html
			}
		}

	}
	return html
}
