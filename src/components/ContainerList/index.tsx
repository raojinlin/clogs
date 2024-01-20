import React from 'react';
import { getContainers } from '../../services/container.ts';
import { Button, Space, Table, Tag } from 'antd';
import moment from 'moment';

const colors: Array<string> = ["magenta", "red", "gold", "orange", "green", "yellow"];

export default function ContainerList() {
    const [containers, setContainers] = React.useState([]);

    React.useEffect(() => {
        getContainers().then((containers) => setContainers(containers));
    }, []);

    const columns = [
        {
            title: 'ID',
            dataIndex: 'id',
            render(id) {
                return id.slice(0, 8);
            }
        },
        {
            title: '名称',
            dataIndex: 'name',
            render(name) {
                return name.slice(1);
            }
        },
        {
            title: '状态',
            dataIndex: 'state',
        },
        {
            title: '创建时间',
            dataIndex: 'created',
            render(created) {
                return moment(created * 1000).format('YYYY-MM-DD HH:mm:ss');
            }
        },
        {
            title: '操作',
            width: 100,
            render(_, container) {
                return (
                    <Button type='link'>查看日志</Button>
                );
            }
        }
    ]

    return (
        <div>
            <Table dataSource={containers} columns={columns} />
        </div>
    )
}