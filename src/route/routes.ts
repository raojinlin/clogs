import Layout from '../Layout/index.tsx';
import ContainerList from '../components/ContainerList/index.tsx';

export const routes = [
  {
    path: '',
    Component: Layout,
    children: [
      {
        path: '/',
        title: 'Container List',
        Component: ContainerList,
      }
    ]
  }
]