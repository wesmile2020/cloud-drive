import styles from './LinkWrapper.module.css';

function LinkWrapper(props: React.PropsWithChildren) {
  return (
    <div className={styles.link_wrapper}>
      {props.children}
    </div>
  )  

}

export default LinkWrapper;
