import styles from './CardWrapper.module.css';

interface Props {
  title: string
}

function CardWrapper(props: React.PropsWithChildren<Props>) {
  return (
    <div className={styles.card_container}>
      <div className={styles.card}>
        <h1 className={styles.title}>{props.title}</h1>
        {props.children}
      </div>
    </div>
  );
}

export default CardWrapper;
