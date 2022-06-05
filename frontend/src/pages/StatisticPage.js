import React from "react";
import {observer} from "mobx-react";
import "./StatisticPage.css"

const StatisticItem = (props) => {
  return (
    <div className={"statistic__item"}>
      <div className={'statistic__item__name'}>
        {props.name}
      </div>
      <div className={'statistic__item__value'}>
        {props.value}
      </div>
    </div>
  )
};

@observer
class StatisticPage extends React.Component {
  constructor(props) {
    super(props);
  }

  componentDidMount() {
    const statisticStore = this.props.statisticStore;
    statisticStore.get()
  };

  render() {
    const statistic = this.props.statisticStore.statistics;

    console.log(statistic);
    return (
      <div className={'statistics-page'}>
        <StatisticItem name={'Зарегестрировано пользователей за сутки:'} value={statistic.users}/>
        <StatisticItem name={'Создано постов за сутки:'} value={statistic.posts}/>
        <StatisticItem name={'Создано веток за сутки:'} value={statistic.threads}/>
        <StatisticItem name={'Оставлено голосов за сутки:'} value={statistic.votes}/>
      </div>
    );
  }
}

export default StatisticPage;