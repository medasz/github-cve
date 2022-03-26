package crawlers

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github-cve/db/db_affected_version"
	"github-cve/db/db_metasploit_module"
	"github-cve/db/db_product_affected"
	"github-cve/db/db_references"
	"github-cve/db/db_vulnerability"
	"github-cve/model/affected_version"
	"github-cve/model/metasploit_module"
	"github-cve/model/product_affected"
	"github-cve/model/references"
	"github-cve/model/vulnerability"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"gorm.io/gorm"
)

type Github struct {
}

func (g *Github) Run() {
	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
		colly.AllowedDomains("www.cvedetails.com"),
	)

	// 爬取页数
	c.OnHTML("#pagingb", func(element *colly.HTMLElement) {
		element.ForEach("a", func(i int, element *colly.HTMLElement) {
			url := element.Attr("href")
			url = fmt.Sprintf("https://www.cvedetails.com/%s", url)
			fmt.Println(url)
			element.Request.Visit(url)
		})
	})

	// 爬取每页的项
	c.OnHTML("body", func(element *colly.HTMLElement) {
		element.DOM.Find("#vulnslisttable tbody tr").Each(func(i int, selection *goquery.Selection) {
			selection.Find("td").Each(func(i int, selection *goquery.Selection) {
				if i == 9 {
					element.Request.Ctx.Put("Access", selection.Text())
				} else if i == 4 {
					fmt.Printf("=======%s======\n", strings.TrimSpace(selection.Text()))
					element.Request.Ctx.Put("VulnerabilityType", strings.TrimSpace(selection.Text()))
				}
			})
			url, exists := selection.Find("a[href]").Attr("href")
			if exists {
				url = fmt.Sprintf("https://www.cvedetails.com%s", url)
				fmt.Println(url)
				element.Request.Visit(url)
			}
		})
	})

	// 爬取漏洞细节
	c.OnHTML("#contentdiv #cvedetails", func(element *colly.HTMLElement) {
		data := new(vulnerability.Vulnerability)
		data.CVEID = g.GetCVEID(element)
		data.VulnerabilityDetails = g.GetVulnerabilityDetails(element)
		data.PublishDate = g.GetPublishDate(element)
		data.LastUpdateDate = g.GetLastUpdateDate(element)
		element.DOM.Find("#cvssscorestable tbody tr").Each(func(i int, selection *goquery.Selection) {
			content := selection.Find("td").Text()
			content = strings.TrimSpace(content)
			switch i {
			case 0:
				data.CVSSScore = content
			case 1:
				res := strings.Split(content, "\n")
				if len(res) < 2 {
					data.ConfidentialityImpact = fmt.Sprintf("%s|%s", res[0], "")
				} else {
					data.ConfidentialityImpact = fmt.Sprintf("%s|%s", res[0], strings.TrimSpace(res[1]))
				}
			case 2:
				res := strings.Split(content, "\n")
				if len(res) < 2 {
					data.IntegrityImpact = fmt.Sprintf("%s|%s", res[0], "")
				} else {
					data.IntegrityImpact = fmt.Sprintf("%s|%s", res[0], strings.TrimSpace(res[1]))
				}
			case 3:
				res := strings.Split(content, "\n")
				if len(res) < 2 {
					data.AvailabilityImpact = fmt.Sprintf("%s|%s", res[0], "")
				} else {
					data.AvailabilityImpact = fmt.Sprintf("%s|%s", res[0], strings.TrimSpace(res[1]))
				}
			case 4:
				res := strings.Split(content, "\n")
				if len(res) < 2 {
					data.AccessComplexity = fmt.Sprintf("%s|%s", res[0], "")
				} else {
					data.AccessComplexity = fmt.Sprintf("%s|%s", res[0], strings.TrimSpace(res[1]))
				}
			case 5:
				res := strings.Split(content, "\n")
				if len(res) < 2 {
					data.Authentication = fmt.Sprintf("%s|%s", res[0], "")
				} else {
					data.Authentication = fmt.Sprintf("%s|%s", res[0], strings.TrimSpace(res[1]))
				}
			case 6:
				data.GainedAccess = content
			case 7:
				data.VulnerabilityType = content
			case 8:
				data.CWEID = content
				val, exists := selection.Find("a").Attr("href")
				if exists {
					data.CWEID += "|" + val
				}
			}
			fmt.Printf("--%s--\n", content)
		})
		data.Access = element.Request.Ctx.Get("Access")
		data.VulnerabilityTypeSimple = element.Request.Ctx.Get("VulnerabilityType")
		//fmt.Printf("============%s==============\n", data.Access)
		//fmt.Printf("============%s==============\n", data.VulnerabilityTypeSimple)
		res, err := db_vulnerability.GetItemByCveId(data.CVEID)
		if err != nil && err != gorm.ErrRecordNotFound {
			log.Println(err)
			return
		}
		var vulnerabilityId uint
		if res.ID != 0 {
			//数据存在
			err := db_vulnerability.FirstOrCreate(data)
			if err != nil {
				log.Println(err)
				return
			}
			if data.ID != res.ID {
				//数据更新
				vulnerabilityId = data.ID
				db_vulnerability.DeleteItemByVulnerabilityId(res.ID)
				db_references.DeleteItemByVulnerabilityId(res.ID)
				db_affected_version.DeleteItemByVulnerabilityId(res.ID)
				db_metasploit_module.DeleteItemByVulnerabilityId(res.ID)
				db_product_affected.DeleteItemByVulnerabilityId(res.ID)
			}
		} else {
			err = db_vulnerability.FirstOrCreate(data)
			if err != nil {
				log.Println(err)
				return
			}
		}
		vulnerabilityId = data.ID

		element.DOM.Find("#vulnprodstable tbody tr").Each(func(i int, selection *goquery.Selection) {
			if i > 0 {
				productAffected := new(product_affected.ProductAffected)
				selection.Find("td").Each(func(i int, selection *goquery.Selection) {
					switch i {
					case 1:
						fmt.Println(strings.TrimSpace(selection.Text()))
						productAffected.ProductType = strings.TrimSpace(selection.Text())
					case 2:
						fmt.Println(strings.TrimSpace(selection.Text()))
						productAffected.Vendor = strings.TrimSpace(selection.Text())
						if tmp, exists := selection.Find("a").Attr("href"); exists {
							productAffected.Vendor += "|https://www.cvedetails.com" + tmp
						}
					case 3:
						fmt.Println(strings.TrimSpace(selection.Text()))
						productAffected.Product = strings.TrimSpace(selection.Text())
						if tmp, exists := selection.Find("a").Attr("href"); exists {
							productAffected.Product += "|https://www.cvedetails.com" + tmp
						}
					case 4:
						productAffected.Version = strings.TrimSpace(selection.Text())
					case 5:
						productAffected.Update = strings.TrimSpace(selection.Text())
					case 6:
						productAffected.Edition = strings.TrimSpace(selection.Text())
					case 7:
						productAffected.Language = strings.TrimSpace(selection.Text())
					case 8:
						selection.Find("a").Each(func(i int, selection *goquery.Selection) {
							switch i {
							case 0:
								href, exists := selection.Attr("href")
								if exists {
									productAffected.VersionDetail = strings.TrimSpace(selection.Text()) + "|https://www.cvedetails.com/" + href
								}
							case 1:
								href, exists := selection.Attr("href")
								if exists {
									productAffected.Vulnerabilities = strings.TrimSpace(selection.Text()) + "|https://www.cvedetails.com/" + href
								}
							}
						})
					}
				})
				productAffected.VulnerabilityID = vulnerabilityId
				db_product_affected.FirstOrCreate(productAffected)
			}
		})

		element.DOM.Find("#vulnversconuttable tbody tr").Each(func(i int, selection *goquery.Selection) {
			if i > 0 {
				affectedVersion := new(affected_version.AffectedVersion)
				selection.Find("td").Each(func(i int, selection *goquery.Selection) {
					switch i {
					case 0:
						href, exists := selection.Find("a").Attr("href")
						if exists {
							affectedVersion.Vendor = strings.TrimSpace(selection.Find("a").Text()) + "|https://www.cvedetails.com/" + href
						}
					case 1:
						href, exists := selection.Find("a").Attr("href")
						if exists {
							affectedVersion.Product = strings.TrimSpace(selection.Find("a").Text()) + "|https://www.cvedetails.com/" + href
						}
					case 2:
						affectedVersion.VulnerableVersions = strings.TrimSpace(selection.Text())
					}
				})
				affectedVersion.VulnerabilityID = vulnerabilityId
				db_affected_version.FirstOrCreate(affectedVersion)
			}
		})

		element.DOM.Find("#vulnrefstable tbody tr").Each(func(i int, selection *goquery.Selection) {
			reference := new(references.Reference)
			fmt.Println("^^^^^^^^^^", selection.Text())
			selection.Find("td").Each(func(i int, selection *goquery.Selection) {
				href, exists := selection.Find("a").Attr("href")
				if exists {
					reference.Link = href
					reference.Content = strings.TrimSpace(strings.ReplaceAll(selection.Text(), href, ""))
				}
			})
			reference.VulnerabilityID = vulnerabilityId
			fmt.Println(reference)
			db_references.FirstOrCreate(reference)
		})
		element.Request.Ctx.Put("vulnerabilityId", vulnerabilityId)
	})

	c.OnHTML("#contentdiv #metasploitmodstable", func(element *colly.HTMLElement) {
		vulnerabilityIdRaw := element.Request.Ctx.GetAny("vulnerabilityId")
		fmt.Println(vulnerabilityIdRaw)
		if vulnerabilityId, ok := vulnerabilityIdRaw.(uint); ok {
			metasploitModule := new(metasploit_module.MetasploitModule)
			metasploitModule.VulnerabilityID = uint(vulnerabilityId)
			metasploitModule.Content = element.Text
			db_metasploit_module.FirstOrCreate(metasploitModule)
		}
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// 开始爬取数据
	for true {
		for i := 1988; i < time.Now().Year(); i++ {
			// Start scraping on https://www.cvedetails.com
			c.Visit(fmt.Sprintf("https://www.cvedetails.com/vulnerability-list/year-%d/vulnerabilities.html", i))
		}
		time.Sleep(time.Hour)
	}
}

func (g *Github) GetCVEID(element *colly.HTMLElement) string {
	return element.DOM.Find("h1 a").Text()
}

func (g *Github) GetVulnerabilityDetails(element *colly.HTMLElement) string {
	tmp := element.DOM.Find("div .cvedetailssummary").Text()
	content := strings.TrimSpace(tmp)
	content = strings.ReplaceAll(content, "\t", "")
	content = strings.ReplaceAll(content, "\n\n", "\n")
	contents := strings.Split(content, "\n")
	return contents[0]
}

func (g *Github) GetPublishDate(element *colly.HTMLElement) time.Time {
	tmp := element.DOM.Find("div .cvedetailssummary").Text()
	content := strings.TrimSpace(tmp)
	content = strings.ReplaceAll(content, "\t", "")
	content = strings.ReplaceAll(content, "\n\n", "\n")
	contents := strings.Split(content, "\n")
	reg := regexp.MustCompile("^Publish Date : (.*?)Last Update Date : (.*?)$")
	res := reg.FindStringSubmatch(contents[1])
	publishDate, _ := time.ParseInLocation("2006-01-02", res[1], time.Local)
	return publishDate
}

func (g *Github) GetLastUpdateDate(element *colly.HTMLElement) time.Time {
	tmp := element.DOM.Find("div .cvedetailssummary").Text()
	content := strings.TrimSpace(tmp)
	content = strings.ReplaceAll(content, "\t", "")
	content = strings.ReplaceAll(content, "\n\n", "\n")
	contents := strings.Split(content, "\n")
	reg := regexp.MustCompile("^Publish Date : (.*?)Last Update Date : (.*?)$")
	res := reg.FindStringSubmatch(contents[1])
	lastUpdateDate, _ := time.ParseInLocation("2006-01-02", res[2], time.Local)
	return lastUpdateDate
}
