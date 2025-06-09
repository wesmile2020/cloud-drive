import { FileZipTwoTone, FileTwoTone, FileWordTwoTone, FilePdfTwoTone, FileTextTwoTone, FileImageTwoTone } from '@ant-design/icons';
import { Tooltip } from 'antd';

interface Props {
  name: string;
}

function getPrefix(name: string) {
  if (/\.zip$/.test(name)) {
    return <FileZipTwoTone />
  } else if (/\.docx$/.test(name)) {
    return <FileWordTwoTone /> 
  } else if (/\.pdf$/.test(name)) {
    return <FilePdfTwoTone />
  } else if (/\.(te?xt)$/.test(name)) {
    return <FileTextTwoTone />  
  } else if (/\.(jpe?g|svg|png|gif)$/.test(name)) {
    return <FileImageTwoTone />
  }
  return <FileTwoTone />
}

function FileName(props: Props) {
  return (
    <Tooltip title={props.name}
      placement='topLeft'>
      {getPrefix(props.name)}
      &nbsp;
      {props.name}
    </Tooltip>
  );
}

export default FileName;
