package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/urfave/cli"
)

const (
	appName    = "nagios-sqs"
	appVersion = "1.0"
)

var (
	awsConfiguration awsConf
)

type awsConf struct {
	Region string
	Key    string
	Secret string
	QURL   string
}

type nagiosService struct {
	NotificationFrom       string `json:"NotificationFrom"`
	Notificationtype       string `json:"Notificationtype"`
	Longdatetime           string `json:"Longdatetime"`
	Hostname               string `json:"Hostname"`
	Hostdisplayname        string `json:"Hostdisplayname"`
	Hostalias              string `json:"Hostalias"`
	Hostaddress            string `json:"Hostaddress"`
	Hoststate              string `json:"Hoststate"`
	Hoststateid            string `json:"Hoststateid"`
	Lasthoststate          string `json:"Lasthoststate"`
	Lasthoststateid        string `json:"Lasthoststateid"`
	Hoststatetype          string `json:"Hoststatetype"`
	Hostattempt            string `json:"Hostattempt"`
	Maxhostattempts        string `json:"Maxhostattempts"`
	Hosteventid            string `json:"Hosteventid"`
	Lasthosteventid        string `json:"Lasthosteventid"`
	Hostproblemid          string `json:"Hostproblemid"`
	Lasthostproblemid      string `json:"Lasthostproblemid"`
	Hostlatency            string `json:"Hostlatency"`
	Hostexecutiontime      string `json:"Hostexecutiontime"`
	Hostduration           string `json:"Hostduration"`
	Hostdurationsec        string `json:"Hostdurationsec"`
	Hostdowntime           string `json:"Hostdowntime"`
	Hostpercentchange      string `json:"Hostpercentchange"`
	Hostgroupname          string `json:"Hostgroupname"`
	Hostgroupnames         string `json:"Hostgroupnames"`
	Lasthostcheck          string `json:"Lasthostcheck"`
	Lasthoststatechange    string `json:"Lasthoststatechange"`
	Lasthostup             string `json:"Lasthostup"`
	Lasthostdown           string `json:"Lasthostdown"`
	Lasthostunreachable    string `json:"Lasthostunreachable"`
	Hostoutput             string `json:"Hostoutput"`
	Longhostoutput         string `json:"Longhostoutput"`
	Hostperfdata           string `json:"Hostperfdata"`
	Servicedesc            string `json:"Servicedesc"`
	Servicedisplayname     string `json:"Servicedisplayname"`
	Servicestate           string `json:"Servicestate"`
	Servicestateid         string `json:"Servicestateid"`
	Lastservicestate       string `json:"Lastservicestate"`
	Lastservicestateid     string `json:"Lastservicestateid"`
	Servicestatetype       string `json:"Servicestatetype"`
	Serviceattempt         string `json:"Serviceattempt"`
	Maxserviceattempts     string `json:"Maxserviceattempts"`
	Serviceisvolatile      string `json:"Serviceisvolatile"`
	Serviceeventid         string `json:"Serviceeventid"`
	Lastserviceeventid     string `json:"Lastserviceeventid"`
	Serviceproblemid       string `json:"Serviceproblemid"`
	Lastserviceproblemid   string `json:"Lastserviceproblemid"`
	Servicelatency         string `json:"Servicelatency"`
	Serviceexecutiontime   string `json:"Serviceexecutiontime"`
	Serviceduration        string `json:"Serviceduration"`
	Servicedurationsec     string `json:"Servicedurationsec"`
	Servicedowntime        string `json:"Servicedowntime"`
	Servicepercentchange   string `json:"Servicepercentchange"`
	Servicegroupname       string `json:"Servicegroupname"`
	Servicegroupnames      string `json:"Servicegroupnames"`
	Lastservicecheck       string `json:"Lastservicecheck"`
	Lastservicestatechange string `json:"Lastservicestatechange"`
	Lastserviceok          string `json:"Lastserviceok"`
	Lastservicewarning     string `json:"Lastservicewarning"`
	Lastserviceunknown     string `json:"Lastserviceunknown"`
	Lastservicecritical    string `json:"Lastservicecritical"`
	Serviceoutput          string `json:"Serviceoutput"`
	Longserviceoutput      string `json:"Longserviceoutput"`
	Serviceperfdata        string `json:"Serviceperfdata"`
	Servicenotesurl        string `json:"Servicenotesurl"`
	Hostnotesurl           string `json:"Hostnotesurl"`
}

var (
	nagiosSvc nagiosService
)

func main() {

	app := cli.NewApp()
	app.Name = appName
	app.HelpName = appName
	app.Usage = "Simple CLI to send Nagios alerts to AWS SQS"
	app.Version = appVersion
	app.Copyright = ""
	app.Authors = []cli.Author{
		{
			Name: "Petr Artamonov",
		},
		{
			Name: "Rafal Pieniazek",
		},
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "aws-accesskey",
			Value:       "",
			Usage:       "AWS Access key",
			Destination: &awsConfiguration.Key,
		},
		cli.StringFlag{
			Name:        "aws-secretkey",
			Value:       "nagios",
			Usage:       "AWS Secret key",
			Destination: &awsConfiguration.Secret,
		},
		cli.StringFlag{
			Name:        "aws-region",
			Value:       "",
			Usage:       "AWS Region",
			Destination: &awsConfiguration.Region,
		},
		cli.StringFlag{
			Name:        "aws-qurl",
			Value:       "",
			Usage:       "SQS URL",
			Destination: &awsConfiguration.QURL,
		},
		cli.StringFlag{
			Name:        "nfr",
			Value:       "",
			Usage:       "NotificationFrom",
			Destination: &nagiosSvc.NotificationFrom,
		},
		cli.StringFlag{
			Name:        "t",
			Value:       "",
			Usage:       "NOTIFICATIONTYPE",
			Destination: &nagiosSvc.Notificationtype,
		},
		cli.StringFlag{
			Name:        "ldt",
			Value:       "",
			Usage:       "LONGDATETIME",
			Destination: &nagiosSvc.Longdatetime,
		},
		cli.StringFlag{
			Name:        "hn",
			Value:       "",
			Usage:       "HOSTNAME",
			Destination: &nagiosSvc.Hostname,
		},
		cli.StringFlag{
			Name:        "hnu",
			Value:       "",
			Usage:       "HOSTNOTESURL",
			Destination: &nagiosSvc.Hostnotesurl,
		},
		cli.StringFlag{
			Name:        "hdn",
			Value:       "",
			Usage:       "HOSTDISPLAYNAME",
			Destination: &nagiosSvc.Hostdisplayname,
		},
		cli.StringFlag{
			Name:        "hal",
			Value:       "",
			Usage:       "HOSTALIAS",
			Destination: &nagiosSvc.Hostalias,
		},
		cli.StringFlag{
			Name:        "haddr",
			Value:       "",
			Usage:       "HOSTADDRESS",
			Destination: &nagiosSvc.Hostaddress,
		},
		cli.StringFlag{
			Name:        "hs",
			Value:       "",
			Usage:       "HOSTSTATE",
			Destination: &nagiosSvc.Hoststate,
		},
		cli.StringFlag{
			Name:        "hsi",
			Value:       "",
			Usage:       "HOSTSTATEID",
			Destination: &nagiosSvc.Hoststateid,
		},
		cli.StringFlag{
			Name:        "lhsi",
			Value:       "",
			Usage:       "LASTHOSTSTATEID",
			Destination: &nagiosSvc.Lasthoststateid,
		},
		cli.StringFlag{
			Name:        "lhs",
			Value:       "",
			Usage:       "LASTHOSTSTATE",
			Destination: &nagiosSvc.Lasthoststate,
		},
		cli.StringFlag{
			Name:        "hst",
			Value:       "",
			Usage:       "HOSTSTATETYPE",
			Destination: &nagiosSvc.Hoststatetype,
		},
		cli.StringFlag{
			Name:        "ha",
			Usage:       "HOSTATTEMPT",
			Destination: &nagiosSvc.Hostattempt,
		},
		cli.StringFlag{
			Name:        "mha",
			Usage:       "MAXHOSTATTEMPTS",
			Destination: &nagiosSvc.Maxhostattempts,
		},
		cli.StringFlag{
			Name:        "hei",
			Usage:       "HOSTEVENTID",
			Destination: &nagiosSvc.Hosteventid,
		},
		cli.StringFlag{
			Name:        "lhei",
			Usage:       "LASTHOSTEVENTID",
			Destination: &nagiosSvc.Lasthosteventid,
		},
		cli.StringFlag{
			Name:        "hpi",
			Value:       "",
			Usage:       "HOSTPROBLEMID",
			Destination: &nagiosSvc.Hostproblemid,
		},
		cli.StringFlag{
			Name:        "lhpi",
			Value:       "",
			Usage:       "LASTHOSTPROBLEMID",
			Destination: &nagiosSvc.Lasthostproblemid,
		},
		cli.StringFlag{
			Name:        "hl",
			Value:       "",
			Usage:       "HOSTLATENCY",
			Destination: &nagiosSvc.Hostlatency,
		},
		cli.StringFlag{
			Name:        "het",
			Value:       "",
			Usage:       "HOSTEXECUTIONTIME",
			Destination: &nagiosSvc.Hostexecutiontime,
		},
		cli.StringFlag{
			Name:        "hd",
			Value:       "",
			Usage:       "HOSTDURATION",
			Destination: &nagiosSvc.Hostduration,
		},
		cli.StringFlag{
			Name:        "hds",
			Value:       "",
			Usage:       "HOSTDURATIONSEC",
			Destination: &nagiosSvc.Hostdurationsec,
		},
		cli.StringFlag{
			Name:        "hdt",
			Value:       "",
			Usage:       "HOSTDOWNTIME",
			Destination: &nagiosSvc.Hostdowntime,
		},
		cli.StringFlag{
			Name:        "hpc",
			Value:       "",
			Usage:       "HOSTPERCENTCHANGE",
			Destination: &nagiosSvc.Hostpercentchange,
		},
		cli.StringFlag{
			Name:        "hgn",
			Value:       "",
			Usage:       "HOSTGROUPNAME",
			Destination: &nagiosSvc.Hostgroupname,
		},
		cli.StringFlag{
			Name:        "hgns",
			Value:       "",
			Usage:       "HOSTGROUPNAMES",
			Destination: &nagiosSvc.Hostgroupnames,
		},
		cli.StringFlag{
			Name:        "lhc",
			Value:       "",
			Usage:       "LASTHOSTCHECK",
			Destination: &nagiosSvc.Lasthostcheck,
		},
		cli.StringFlag{
			Name:        "lhsc",
			Value:       "",
			Usage:       "LASTHOSTSTATECHANGE",
			Destination: &nagiosSvc.Lasthoststatechange,
		},
		cli.StringFlag{
			Name:        "lhu",
			Value:       "",
			Usage:       "LASTHOSTUP",
			Destination: &nagiosSvc.Lasthostup,
		},
		cli.StringFlag{
			Name:        "lhd",
			Value:       "",
			Usage:       "LASTHOSTDOWN",
			Destination: &nagiosSvc.Lasthostdown,
		},
		cli.StringFlag{
			Name:        "lhur",
			Value:       "",
			Usage:       "LASTHOSTUNREACHABLE",
			Destination: &nagiosSvc.Lasthostunreachable,
		},
		cli.StringFlag{
			Name:        "ho",
			Value:       "",
			Usage:       "HOSTOUTPUT",
			Destination: &nagiosSvc.Hostoutput,
		},
		cli.StringFlag{
			Name:        "lho",
			Value:       "",
			Usage:       "LONGHOSTOUTPUT",
			Destination: &nagiosSvc.Longhostoutput,
		},
		cli.StringFlag{
			Name:        "hpd",
			Value:       "",
			Usage:       "HOSTPERFDATA",
			Destination: &nagiosSvc.Hostperfdata,
		},
		cli.StringFlag{
			Name:        "s",
			Value:       "",
			Usage:       "SERVICEDESC",
			Destination: &nagiosSvc.Servicedesc,
		},
		cli.StringFlag{
			Name:        "sdn",
			Value:       "",
			Usage:       "SERVICEDISPLAYNAME",
			Destination: &nagiosSvc.Servicedisplayname,
		},
		cli.StringFlag{
			Name:        "ss",
			Value:       "",
			Usage:       "SERVICESTATE",
			Destination: &nagiosSvc.Servicestate,
		},
		cli.StringFlag{
			Name:        "ssi",
			Value:       "",
			Usage:       "SERVICESTATEID",
			Destination: &nagiosSvc.Servicestateid,
		},
		cli.StringFlag{
			Name:        "lss",
			Value:       "",
			Usage:       "LASTSERVICESTATE",
			Destination: &nagiosSvc.Lastservicestate,
		},
		cli.StringFlag{
			Name:        "lssi",
			Value:       "",
			Usage:       "LASTSERVICESTATEID",
			Destination: &nagiosSvc.Lastservicestateid,
		},
		cli.StringFlag{
			Name:        "sst",
			Value:       "",
			Usage:       "SERVICESTATETYPE",
			Destination: &nagiosSvc.Servicestatetype,
		},
		cli.StringFlag{
			Name:        "sa",
			Value:       "",
			Usage:       "SERVICEATTEMPT",
			Destination: &nagiosSvc.Serviceattempt,
		},
		cli.StringFlag{
			Name:        "msa",
			Value:       "",
			Usage:       "MAXSERVICEATTEMPTS",
			Destination: &nagiosSvc.Maxserviceattempts,
		},
		cli.StringFlag{
			Name:        "siv",
			Value:       "",
			Usage:       "SERVICEISVOLATILE",
			Destination: &nagiosSvc.Serviceisvolatile,
		},
		cli.StringFlag{
			Name:        "sei",
			Value:       "",
			Usage:       "SERVICEEVENTID",
			Destination: &nagiosSvc.Serviceeventid,
		},
		cli.StringFlag{
			Name:        "lsei",
			Value:       "",
			Usage:       "LASTSERVICEEVENTID",
			Destination: &nagiosSvc.Lastserviceeventid,
		},
		cli.StringFlag{
			Name:        "spi",
			Value:       "",
			Usage:       "SERVICEPROBLEMID",
			Destination: &nagiosSvc.Serviceproblemid,
		},
		cli.StringFlag{
			Name:        "lspi",
			Value:       "",
			Usage:       "LASTSERVICEPROBLEMID",
			Destination: &nagiosSvc.Lastserviceproblemid,
		},
		cli.StringFlag{
			Name:        "sl",
			Value:       "",
			Usage:       "SERVICELATENCY",
			Destination: &nagiosSvc.Servicelatency,
		},
		cli.StringFlag{
			Name:        "set",
			Value:       "",
			Usage:       "SERVICEEXECUTIONTIME",
			Destination: &nagiosSvc.Serviceexecutiontime,
		},
		cli.StringFlag{
			Name:        "sd",
			Value:       "",
			Usage:       "SERVICEDURATION",
			Destination: &nagiosSvc.Serviceduration,
		},
		cli.StringFlag{
			Name:        "sds",
			Value:       "",
			Usage:       "SERVICEDURATIONSEC",
			Destination: &nagiosSvc.Servicedurationsec,
		},
		cli.StringFlag{
			Name:        "sdt",
			Value:       "",
			Usage:       "SERVICEDOWNTIME",
			Destination: &nagiosSvc.Servicedowntime,
		},
		cli.StringFlag{
			Name:        "spc",
			Value:       "",
			Usage:       "SERVICEPERCENTCHANGE",
			Destination: &nagiosSvc.Servicepercentchange,
		},
		cli.StringFlag{
			Name:        "sgn",
			Value:       "",
			Usage:       "SERVICEGROUPNAME",
			Destination: &nagiosSvc.Servicegroupname,
		},
		cli.StringFlag{
			Name:        "sgns",
			Value:       "",
			Usage:       "SERVICEGROUPNAMES",
			Destination: &nagiosSvc.Servicegroupnames,
		},
		cli.StringFlag{
			Name:        "lsch",
			Value:       "",
			Usage:       "LASTSERVICECHECK",
			Destination: &nagiosSvc.Lastservicecheck,
		},
		cli.StringFlag{
			Name:        "lssc",
			Value:       "",
			Usage:       "LASTSERVICESTATECHANGE",
			Destination: &nagiosSvc.Lastservicestatechange,
		},
		cli.StringFlag{
			Name:        "lsok",
			Value:       "",
			Usage:       "LASTSERVICEOK",
			Destination: &nagiosSvc.Lastserviceok,
		},
		cli.StringFlag{
			Name:        "lsw",
			Value:       "",
			Usage:       "LASTSERVICEWARNING",
			Destination: &nagiosSvc.Lastservicewarning,
		},
		cli.StringFlag{
			Name:        "lsu",
			Value:       "",
			Usage:       "LASTSERVICEUNKNOWN",
			Destination: &nagiosSvc.Lastserviceunknown,
		},
		cli.StringFlag{
			Name:        "lsc",
			Value:       "",
			Usage:       "LASTSERVICECRITICAL",
			Destination: &nagiosSvc.Lastservicecritical,
		},
		cli.StringFlag{
			Name:        "so",
			Value:       "",
			Usage:       "SERVICEOUTPUT",
			Destination: &nagiosSvc.Serviceoutput,
		},
		cli.StringFlag{
			Name:        "lso",
			Value:       "",
			Usage:       "LONGSERVICEOUTPUT",
			Destination: &nagiosSvc.Longserviceoutput,
		},
		cli.StringFlag{
			Name:        "spd",
			Value:       "",
			Usage:       "SERVICEPERFDATA",
			Destination: &nagiosSvc.Serviceperfdata,
		},
		cli.StringFlag{
			Name:        "snu",
			Value:       "",
			Usage:       "SERVICENOTESURL",
			Destination: &nagiosSvc.Servicenotesurl,
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	app.Action = func(c *cli.Context) error {

		sendMessage(c)

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}

}

func sendMessage(c *cli.Context) {

	ses, err := session.NewSession(&aws.Config{
		Region:      aws.String(awsConfiguration.Region),
		Credentials: credentials.NewStaticCredentials(awsConfiguration.Key, awsConfiguration.Secret, ""),
	})
	if err != nil {
		fmt.Println("Error ses", err)
		return
	}

	// Create a SQS client from just a session.
	svc := sqs.New(ses, &aws.Config{
		Region: aws.String(awsConfiguration.Region),
	})

	// URL of our queue
	qURL := awsConfiguration.QURL

	b, err := json.Marshal(nagiosSvc)

	if err != nil {
		fmt.Println("Error marshall", err)
		return
	}

	result, err := svc.SendMessage(&sqs.SendMessageInput{
		// MessageDeduplicationId: aws.String("xxx"),
		// MessageGroupId:         aws.String("nagios"),
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			"source": &sqs.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String("nagios"), // Here we comply with rule that every event needs to have source
			},
		},
		MessageBody: aws.String(string(b)),
		QueueUrl:    &qURL,
	})

	if err != nil {
		fmt.Println("Error send message", err)
		return
	}

	fmt.Println(*result)
	fmt.Println("Success", *result.MessageId)
}
