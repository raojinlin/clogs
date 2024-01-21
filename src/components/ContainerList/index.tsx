import React from 'react';
import { getContainers } from '../../services/container.ts';
import {Button, Modal, Table} from 'antd';
import moment from 'moment';
import './index.css';
import LogScreen from "../LogScreen/index.tsx";

const statusName = {
    running: '正在运行',
    exited: '已退出',
    created: '已创建',
}

export default function ContainerList() {
    const [containers, setContainers] = React.useState([]);
    const [currentContainer, setCurrentContainer] = React.useState(null);

    React.useEffect(() => {
        getContainers().then((containers) =>
            setContainers(containers)
        );
    }, []);

    const columns = [
        {
            title: 'ID',
            dataIndex: 'id',
            key: 'id',
            render(id) {
                return id.slice(0, 8);
            }
        },
        {
            title: '名称',
            dataIndex: 'name',
            key: 'name',
            width: 100,
            render(name) {
                return <div className={'container-name'}>{name.slice(1)}</div>;
            }
        },
        {
            title: '状态',
            dataIndex: 'state',
            render(status: string, container) {
                if (status === 'created') {
                    return statusName[status];
                }
                return <span>{statusName[status] || status}, {container.status}</span>;
            }
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
                    <Button type='link' onClick={() => setCurrentContainer(container)}>
                        <span className={'material-symbols-rounded icon'}>feed</span>
                        查看日志
                    </Button>
                );
            }
        }
    ]

    return (
        <div>
            <Table
                pagination={{total: containers.length, pageSize: 5, showSizeChanger: false}}
                dataSource={containers}
                columns={columns}
            />
            {currentContainer ? (
                <Modal
                    size='big'
                    title={`${currentContainer.name.replace(/^\//, '')}的日志`}
                    open
                    onCancel={() => setCurrentContainer(null)}
                    className='log-screen-modal'
                    footer={[
                        <Button type='primary' onClick={() => setCurrentContainer(null)}>取消</Button>
                    ]}
                >
                    <LogScreen containerID={currentContainer.id} />
                </Modal>
            ) : null}
        </div>
    )
}