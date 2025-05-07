import React from 'react';
import { useParams } from 'react-router';

import HomeMenu from '@/components/HomeMenu';
import { getFiles } from '@/services/api';

function Home() {
  const params = useParams();

  React.useEffect(() => {
    console.log(params.id);
    if (params.id) {
      getFiles(params.id);
    }

  }, [params.id]);

  return (
    <>
      <HomeMenu directoryId={Number(params.id)}/>
    </> 
  );
}

export default Home;
