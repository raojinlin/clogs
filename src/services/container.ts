declare type Containers = {
    name: string,
    id: string,
    created: number,
    labels: string,
    status: string,
}


export async function getContainers(): Promise<Containers> {
    return fetch('http://localhost:8082/api/container/list').then(r => r.json());
}