import React from 'react'
import './index.css'
import * as querystring from "querystring";

export default function LogScreen({ containerID, tail=100, logFile='stdout', showStdout=true, showStderr=true }) {
    const [logs, setLogs] = React.useState([]);
    const ele = React.useRef();

    React.useEffect(() => {
        const eventSourceUrl = `http://localhost:8082/api/container/logs/${containerID}?` + querystring.encode({
            logFile,
            showStdout,
            showStderr,
            tail,
        });
        const eventSource = new EventSource(eventSourceUrl);

        eventSource.onmessage = e => {
            const textContent = atob(e.data.replace(/^"|"$/g, ''));
            setLogs((arr) => {
                return [...arr, textContent];
            });
            if (ele.current) {
                ele.current.scrollTo(0, ele.current.offsetHeight);
            }
            hljs.highlightAll();
        };

        eventSource.onerror = () => {
            eventSource.close();
        };

        return () => {
            eventSource.close();
        }
    }, []);

    return (
        <pre ref={ele}>
            {logs.map((line, i) => <code className='hljs language-accesslog' key={i}>{line}</code>)}
        </pre>
    );
}