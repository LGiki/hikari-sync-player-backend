package cowtransfer_parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"hikari_sync_player/entity"
	"hikari_sync_player/pkg/request"
	"io/ioutil"
	"regexp"
)

type TransferDetail struct {
	AllowDownloadAll               bool          `json:"allowDownloadAll"`
	Blocked                        bool          `json:"blocked"`
	BrowseCount                    int           `json:"browseCount"`
	CanOnlyDownloadInMySpace       bool          `json:"canOnlyDownloadInMySpace"`
	CreamSuggestPrice              float32       `json:"creamSuggestPrice"`
	CreamTransfer                  bool          `json:"creamTransfer"`
	CreamTransferPrice             float32       `json:"creamTransferPrice"`
	CreamTransferRequestStatus     interface{}   `json:"creamTransferRequestStatus"`
	CreamTransferRevenue           float32       `json:"creamTransferRevenue"`
	Day                            int           `json:"day"`
	DaysLeft                       int           `json:"daysLeft"`
	DaysToCleanUp                  int           `json:"daysToCleanUp"`
	Deleted                        bool          `json:"deleted"`
	DisabledSaleTransfer           bool          `json:"diabledSaleTransfer"`
	Disabled                       bool          `json:"disabled"`
	DisplayEnterpriseTransferOwner bool          `json:"displayEnterpriseTransferOwner"`
	DownloadLeft                   int           `json:"downloadLeft"`
	DownloadLimit                  int           `json:"downloadLimit"`
	Downloaded                     int           `json:"downloaded"`
	EmailReceivers                 []interface{} `json:"emailReceivers"`
	EnableDownload                 bool          `json:"enableDownload"`
	EnablePreview                  bool          `json:"enablePreview"`
	EnableSingleFileDownload       bool          `json:"enableSingleFileDownload"`
	ErrorCode                      int           `json:"errorCode"`
	ExceedCustomDownloadLimit      bool          `json:"exceeedCustomDownloadLimit"`
	ExpireAt                       string        `json:"expireAt"`
	ExpireHours                    int           `json:"expireHours"`
	FilesAmount                    int           `json:"filesAmount"`
	FirstFile                      FirstFile     `json:"firstFile"`
	FolderFilesAmount              int           `json:"folderFilesAmount"`
	FrontendShow                   bool          `json:"frontendShow"`
	GUID                           string        `json:"guid"`
	HasFolderStructure             bool          `json:"hasFolderStrucutre"`
	IUploaded                      bool          `json:"iuploaded"`
	Language                       string        `json:"language"`
	Message                        string        `json:"message"`
	MiniAppQrCode                  interface{}   `json:"miniAppQrCode"`
	Month                          int           `json:"month"`
	NeedPassword                   bool          `json:"needPassword"`
	Page                           int           `json:"page"`
	Password                       interface{}   `json:"password"`
	ProAccount                     bool          `json:"proAccount"`
	ReceivedTransfer               bool          `json:"receivedTransfer"`
	Receivers                      []interface{} `json:"receivers"`
	Restored                       int           `json:"restored"`
	ShareShow                      bool          `json:"shareShow"`
	SyncedFolderGUID               interface{}   `json:"syncedFolderGuid"`
	TempCode                       string        `json:"tempCode"`
	TempCodeValid                  bool          `json:"tempCodeValid"`
	TotalFileAmount                int           `json:"totalFileAmount"`
	TotalPages                     int           `json:"totalPages"`
	TotalSizeByte                  float64       `json:"totalSizeByte"`
	TotalSizeGb                    float64       `json:"totalSizeGb"`
	TransferFileDtos               []interface{} `json:"transferFileDtos"`
	TransferName                   string        `json:"transferName"`
	UniqueURL                      string        `json:"uniqueUrl"`
	UploadDate                     string        `json:"uploadDate"`
	Uploaded                       bool          `json:"uploaded"`
	ValidDaysLimit                 int64         `json:"validDaysLimit"`
	VersionTag                     int           `json:"versionTag"`
	Year                           int           `json:"year"`
	ZipComplete                    bool          `json:"zipComplete"`
	ZipDownloadFile                interface{}   `json:"zipDownloadFile"`
}

type FirstFile struct {
	AbleToLoadThumbnail      bool        `json:"ableToLoadThumbnail"`
	CanOnlyDownloadInMySpace bool        `json:"canOnlyDownloadInMySpace"`
	ContentType              string      `json:"contentType"`
	CreatedAt                string      `json:"createdAt"`
	DownloadName             string      `json:"downloadName"`
	ErrorCode                int         `json:"errorCode"`
	FileName                 string      `json:"fileName"`
	Folder                   bool        `json:"folder"`
	Gif                      bool        `json:"gif"`
	GUID                     string      `json:"guid"`
	Hash                     string      `json:"hash"`
	Image                    bool        `json:"image"`
	Name                     string      `json:"name"`
	PreviewableImage         bool        `json:"previewableImage"`
	Size                     string      `json:"size"`
	SizeInByte               float32     `json:"sizeInByte"`
	SizeInMB                 string      `json:"sizeInMB"`
	SizeInMBInteger          int         `json:"sizeInMBInteger"`
	TransferGUID             interface{} `json:"transferGuid"`
	Type                     interface{} `json:"type"`
	Uploaded                 bool        `json:"uploaded"`
	URL                      interface{} `json:"url"`
	Valid                    bool        `json:"valid"`
	Video                    bool        `json:"video"`
}

type DownloadLink struct {
	Link string `json:"link"`
	Name string `json:"name"`
}

func getCowTransferFileId(url string) (string, error) {
	cowTransferShareIdPattern := regexp.MustCompile(`(?m)http[s]*://cowtransfer\.com/s/(?P<id>[a-zA-Z0-9]+)/*`)
	cowTransferShareIdMatchResult := cowTransferShareIdPattern.FindStringSubmatch(url)
	if len(cowTransferShareIdMatchResult) == 0 {
		return "", errors.New("fail to parse cow transfer share link")
	}
	return cowTransferShareIdMatchResult[cowTransferShareIdPattern.SubexpIndex("id")], nil
}

func ParseShareLink(url string) (*entity.Video, error) {
	fileId, err := getCowTransferFileId(url)
	if err != nil {
		return nil, err
	}
	transferDetailResponse, err := request.Get(fmt.Sprintf("https://cowtransfer.com/api/transfer/v2/transferdetail?url=%s", fileId))
	if err != nil {
		return nil, err
	}
	defer transferDetailResponse.Body.Close()
	transferDetailResponseBytes, err := ioutil.ReadAll(transferDetailResponse.Body)
	if err != nil {
		return nil, err
	}
	var transferDetail TransferDetail
	err = json.Unmarshal(transferDetailResponseBytes, &transferDetail)
	if err != nil {
		return nil, err
	}
	fileGuid := transferDetail.GUID
	downloadLinkResponse, err := request.Get(fmt.Sprintf("https://cowtransfer.com/api/transfer/all_download_links?guid=%s", fileGuid))
	if err != nil {
		return nil, err
	}
	defer downloadLinkResponse.Body.Close()
	downloadLinkResponseBytes, err := ioutil.ReadAll(downloadLinkResponse.Body)
	if err != nil {
		return nil, err
	}
	var downloadLinks []DownloadLink
	err = json.Unmarshal(downloadLinkResponseBytes, &downloadLinks)
	if err != nil {
		return nil, err
	}
	if len(downloadLinks) == 0 {
		return nil, errors.New("fail to parse cow transfer download link")
	}
	return &entity.Video{
		Url: downloadLinks[0].Link,
	}, nil
}
