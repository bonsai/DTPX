package main

import ("bufio"; "encoding/json"; "flag"; "fmt"; "os"; "strings")
type Record map[string]any
func main(){ in:=flag.String("in","","input interview JSONL"); out:=flag.String("out","toiletpaper.svg","output SVG"); flag.Parse(); if *in==""{fmt.Fprintln(os.Stderr,"-in is required");os.Exit(2)}; rs,e:=readJSONL(*in);if e!=nil{panic(e)};if e=writeSVG(*out,collect(rs));e!=nil{panic(e)};fmt.Println(*out)}
func readJSONL(p string)([]Record,error){f,e:=os.Open(p);if e!=nil{return nil,e};defer f.Close();var out []Record;s:=bufio.NewScanner(f);s.Buffer(make([]byte,1024),1024*1024);for s.Scan(){if strings.TrimSpace(s.Text())==""{continue};var r Record;if e:=json.Unmarshal(s.Bytes(),&r);e!=nil{return nil,e};out=append(out,r)};return out,s.Err()}
func collect(rs []Record)string{var b strings.Builder;for _,r:=range rs{switch r["type"]{case "question":b.WriteString("Q. ");b.WriteString(fmt.Sprint(r["text"]));b.WriteString("\n");case "answer":b.WriteString("A. ");b.WriteString(fmt.Sprint(r["text"]));b.WriteString("\n");case "section":b.WriteString(fmt.Sprint(r["text"]));b.WriteString("\n")}};return b.String()}
func writeSVG(p,t string)error{f,e:=os.Create(p);if e!=nil{return e};defer f.Close();fmt.Fprintln(f,"<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"75mm\" height=\"860mm\" viewBox=\"0 0 75 860\">");fmt.Fprintln(f,"<style>text{font-family:sans-serif;font-size:3.175px}</style>");fmt.Fprintln(f,"<text x=\"5\" y=\"8\">");for _,line:=range strings.Split(t,"\n"){fmt.Fprintf(f,"<tspan x=\"5\" dy=\"4.94\">%s</tspan>",escape(line))};fmt.Fprintln(f,"</text></svg>");return nil}
func escape(s string)string{s=strings.ReplaceAll(s,"&","&amp;");s=strings.ReplaceAll(s,"<","&lt;");return strings.ReplaceAll(s,">","&gt;")}
