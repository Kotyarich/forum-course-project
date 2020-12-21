import React from 'react'
import Button from "./Button";
import './Rating.css'

const Rating = (props) => {
  let className = '';
  if (!props.votes) {
    className = 'zero'
  } else if (props.votes > 0) {
    className = 'positive';
  } else {
    className = 'negative';
  }
  console.log(props.user);
  console.log(!props.user);
  return(
    <div className={'rating'}>
      <div className={'rating__votes rating_' + className}>{props.votes}</div>
      <Button name={'rating-up'}
              action={() => {
                props.onClick(props.id, 1)
              }}
              disabled={!props.user}
              title={'+'}/>
      <Button name={'rating-down'}
              action={() => {
                props.onClick(props.id, -1)
              }}
              disabled={!props.user}
              title={'-'}/>
    </div>
  );
};

export default Rating;