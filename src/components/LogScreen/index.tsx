import React from 'react'
import './index.css'
import * as querystring from "querystring";

export default function LogScreen({ containerID, tail=100, logFile='stdout', showStdout=true, showStderr=true }) {
    const [logs, setLogs] = React.useState([]);
    const ele = React.useRef();

    React.useEffect(() => {
        const eventSourceUrl = `/api/container/logs/${containerID}?` + querystring.encode({
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
        };

        eventSource.onerror = () => {
            eventSource.close();
        };

        return () => {
            eventSource.close();
        }
    }, [containerID, logFile, tail, showStderr, showStdout]);

    React.useEffect(() => {
        // if (!ele.current) return;
        //
        // const observer = new MutationObserver(function (mutations, observer) {
        //     ele.current.scrollTop = `${ele.current.offsetHeight}px`;
        //     hljs.highlightAll();
        // });
        //
        // observer.observe(ele.current, {subtree: true, childList: true})
        //
        // return () => {
        //     observer.disconnect();
        // }
    }, []);

    return (
        <pre ref={ele}>
            {logs.map((line, i) => <code className='hljs language-accesslog' key={i}>{line}</code>)}
        </pre>
    );
}