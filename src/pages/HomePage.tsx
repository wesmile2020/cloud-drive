import HomeMenu from '@/components/HomeMenu';
import { useParams } from 'react-router';

function Home() {
  const params = useParams();
  return (
    <>
      <HomeMenu directoryId={Number(params.id)}/>
    </> 
  );
}

export default Home;
