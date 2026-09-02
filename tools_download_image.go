package email

import (
	"fmt"
	"hash/fnv"
	"strings"
	"sync"

	"github.com/monopolly/file"
	"github.com/monopolly/images"
)

func hashlink(link string) string {
	p := fnv.New32()
	p.Write([]byte(link))
	return fmt.Sprintf("%x", p.Sum(nil))
}

// download all valid images from all links
func downloadImagesOld(links []string, options *Options) (images map[string]*File) {
	options.Log("download images")
	images = make(map[string]*File)
	if !options.Images {
		return
	}
	for _, link := range links {
		link = strings.TrimSpace(link)
		img := downloadImage(link, options)
		if img == nil {
			continue
		}
		images[link] = img
	}
	return
}

// download all valid images from all links
func downloadImages(links []string, options *Options) (images map[string]*File) {
	options.Log("download images")
	images = make(map[string]*File)
	var lock sync.Mutex
	if !options.Images {
		return
	}
	var wg sync.WaitGroup
	wg.Add(len(links))

	for _, link := range links {
		go func(link string, options *Options) {
			defer wg.Done()
			link = strings.TrimSpace(link)
			img := downloadImage(link, options)
			if img == nil {
				return
			}
			img.Link = link
			lock.Lock()
			images[link] = img
			options.Log("image count", len(images))
			lock.Unlock()

		}(link, options)
	}

	wg.Wait()
	options.Log("total images count", len(images))
	return
}

// download, check if valid image
func downloadImage(link string, options *Options) (res *File) {
	options.Log("download image", link)

	//_, temp, err := fasthttp.GetTimeout(nil, link, time.Second*10) //
	temp, err := file.Get(link, options.Proxy)
	if err != nil {
		options.Log(err)
		return
	}
	if temp == nil {
		options.Log("nil image file")
		return
	}

	if options.MinImageSizeBytes > 0 && len(temp) < options.MinImageSizeBytes {
		options.Log("min file len", len(temp))
		return
	}

	options.Log("check image")
	img, er := images.New(temp)
	if er != nil {
		options.Log("non image")
		return
	}

	if options.NoGIF && img.Type == images.TypeGIF {
		options.Log("gif is not allow")
		return
	}

	options.Log("image", img.Ext, img.Width, "x", img.Height, len(temp)/1000, "kb")
	n := new(File)
	n.ID = fmt.Sprintf("%s.%s", hashlink(link), img.Ext)
	n.Body = temp
	n.Size = len(temp)
	n.Image = true
	n.Name = fmt.Sprintf("%s.%s", hashlink(link), img.Ext)
	options.Log("image name", n.Name)
	n.Type = fmt.Sprintf("image/%s", img.Ext)
	options.Log("is image")
	return n
}
