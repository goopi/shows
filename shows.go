package shows

import (
    "encoding/xml"
    "fmt"
    "io"
    "io/ioutil"
    "net/http"
    "net/url"
    "os"
    "os/exec"
    "path"
)

const (
    tvdbSearchUrl = "http://www.thetvdb.com/api/GetSeries.php?seriesname=%s&language=en"
    tvdbShowZipUrl = "http://thetvdb.com/api/%s/series/%s/all/en.zip"
    dataDir = "data/shows/%s"
    maxShows = 10
)

type Shows struct {
    Shows []Show `xml:"Series"`
}

type Show struct {
    Id string `xml:"id"`
    Title string `xml:"SeriesName"`
    Overview string
    Actors string
    Genre string
    FirstAired string
    AirsDayOfWeek string `xml:"Airs_DayOfWeek"`
    AirsTime string `xml:"Airs_Time"`
    Runtime string
    Network string
    Status string
    Banner string `xml:"banner"`
    Poster string `xml:"poster"`
    IMDBId string `xml:"IMDB_ID"`
    LastUpdated string `xml:"lastupdated"`
}

type Episodes struct {
    Show Show `xml:"Series"`
    EpisodeList []Episode `xml:"Episode"`
}

type Episode struct {
    Id string `xml:"id"`
    Season int `xml:"SeasonNumber"`
    Number int `xml:"EpisodeNumber"`
    AbsoluteNumber string `xml:"absolute_number"`
    Name string `xml:"EpisodeName"`
    Overview string
    Director string
    Writer string
    FirstAired string
    Thumb string `xml:"filename"`
    ThumbHeight string `xml:"thumb_height"`
    ThumbWidth string `xml:"thumb_width"`
    IMDBId string `xml:"IMDB_ID"`
    LastUpdated string `xml:"lastupdated"`
}

func (s Show) String() string {
    return fmt.Sprintf("%s (%s)", s.Title, s.Id)
}

func (e Episode) String() string {
    return fmt.Sprintf("S%02dE%02d - %s - %s", e.Season, e.Number, e.Name, e.FirstAired)
}

func Search(q string) []Show {
    url := fmt.Sprintf(tvdbSearchUrl, url.QueryEscape(q))
    response, err := http.Get(url)
    if err != nil { panic(err) }
    defer response.Body.Close()

    var shows []Show
    body, err := ioutil.ReadAll(response.Body)

    if err == nil {
        var s Shows
        xml.Unmarshal(body, &s)
        shows = s.Shows
    }

    if len(shows) > maxShows {
        shows = shows[:maxShows]
    }

    return shows
}

func GetShowFiles(id string) error {
    root, _ := os.Getwd()
    filedest := path.Join(root, fmt.Sprintf(dataDir, id))
    filename := path.Join(filedest, "en.xml")
    os.MkdirAll(filedest, 0777)

    key, err := ioutil.ReadFile("key")
    if err != nil { return err }
    key = key[:len(key) - 1]

    if _, err := os.Stat(filename); os.IsNotExist(err) {
        url := fmt.Sprintf(tvdbShowZipUrl, string(key), url.QueryEscape(id))
        response, err := http.Get(url)
        if err != nil { return err }
        defer response.Body.Close()

        f, err := ioutil.TempFile("", "show_")
        if err != nil { return err }
        defer f.Close()

        io.Copy(f, response.Body)

        cmd := exec.Command("unzip", f.Name(), "-d", filedest)
        err = cmd.Run()
        if err != nil { return err }
    }

    return nil
}

func GetEpisodes(id string) ([]Episode, error) {
    root, _ := os.Getwd()
    filedest := path.Join(root, fmt.Sprintf(dataDir, id))
    filename := path.Join(filedest, "en.xml")

    if _, err := os.Stat(filename); os.IsNotExist(err) {
        err = GetShowFiles(id)
        if err != nil { return nil, err }
    }

    f, err := os.Open(filename)
    if err != nil { return nil, err }
    defer f.Close()
    b, _ := ioutil.ReadAll(f)

    var e Episodes
    xml.Unmarshal(b, &e)

    return e.EpisodeList, nil
}
