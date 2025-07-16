import { useEffect, useState } from "react";

type LinkRes = {
  id: number;
  htmlVersion: string;
  pageTitle: string;
  countOfH1: number;
  countOfH2: number;
  countOfH3: number;
  countOfH4: number;
  countOfH5: number;
  countOfH6: number;
  internalLinksCount: number;
  externalLinksCount: number;
  inaccessibleLinksCount: number;
  hasLoginForm: boolean;
  url: string;
  status: string;
};

export default function App() {
  const [results, setResults] = useState<LinkRes[]>([]);
  // const [current, setCurrent] = useState(0)
  const [userInput, setUserInput] = useState(
    "https://www.scrapingcourse.com/ecommerce/"
  );
  const [q, setQ] = useState<string[]>([]);
  type HeadingKey =
    | "countOfH1"
    | "countOfH2"
    | "countOfH3"
    | "countOfH4"
    | "countOfH5"
    | "countOfH6";
  const [headings, setHeadings] = useState<HeadingKey>("countOfH1");

  useEffect(() => {
    // const eventSource = new EventSource("http://localhost:8080/api/results");

    // eventSource.addEventListener("savedResults", (event: MessageEvent) => {
    //   const data = JSON.parse(event.data);
    //   console.log(data);
    //   setResults(data);
    // });

    // eventSource.addEventListener("result", (event: MessageEvent) => {
    //   const data = JSON.parse(event.data);
    //   console.log(data);
    //   setCurrent(data.ID);
    // });

    // eventSource.onerror = () => {
    //   eventSource.close();
    // };

    // return () => {
    //   eventSource.close();
    // };
    fetch("http://localhost:8080/api/results").then(res=>res.json()).then(data=>setResults(data))
  }, [q]);
  return (
    <main>
      App
      <form
        onSubmit={(e) => {
          e.preventDefault();

          setQ((prev) => [...prev, userInput]);
          setUserInput("")
          fetch("http://localhost:8080/api/add", {
            headers: { ContentType: "application/json" },
            method: "POST",
            body: JSON.stringify({ userInput }),
          }).then((res) => res.json().then((data) => console.log(data)));
        }}
      >
        <input
          type="text"
          value={userInput}
          onChange={(e) => setUserInput(e.target.value)}
        />
        <button type="submit">crawl</button>
      </form>
      <table>
        <thead>
          <tr>
            <th> URL </th>
            <th> Page Title </th>
            <th> HTML Version </th>
            <th>
              Headings
              <select
                name="headings"
                value={headings}
                onChange={(e) => setHeadings(e.target.value as HeadingKey)}
              >
                <option value="countOfH1">H1</option>
                <option value="countOfH2">H2</option>
                <option value="countOfH3">H3</option>
                <option value="countOfH4">H4</option>
                <option value="countOfH5">H5</option>
                <option value="countOfH6">H6</option>
              </select>
            </th>
            <th>Internal Links</th>
            <th>External Links</th>
            <th>Inaccessible Links</th>
            <th>Login form?</th>
            <th>Status</th>
          </tr>
        </thead>
        <tbody>
          {results.map((x, i) => (
            <tr key={i}>
              <td title={x.url}>
                <p className="wrap">{x.url}</p>
              </td>
              <td>
                <p className="wrap">{x.pageTitle}</p>
              </td>
              <td> {x.htmlVersion} </td>
              <td> {x[headings]} </td>
              <td> {x.internalLinksCount} </td>
              <td> {x.externalLinksCount} </td>
              <td> {x.inaccessibleLinksCount} </td>
              <td> {x.hasLoginForm ? "✅" : "❌"} </td>
              <td> {x.url == q.at(-1) ? "Crawling...": x.status}  </td>
            </tr>
          ))}
        </tbody>
      </table>
    </main>
  );
}
