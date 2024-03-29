import {action, toJS} from 'mobx'
import Validator from 'validatorjs';

class FormStore {
  getFlattenedValues = (valueKey = 'value') => {
    let data = {};
    let form = toJS(this.form).fields;
    Object.keys(form).forEach(key => {
      data[key] = form[key][valueKey]
    });
    return data
  };

  @action
  onFieldChange = (field, value) => {
    this.form.fields[field].value = value;
    const validation = new Validator(
      this.getFlattenedValues('value'),
      this.getFlattenedValues('rule'));
    this.form.meta.isValid = validation.passes(validation);
    this.form.fields[field].error = validation.errors.first(field)
  };

  @action
  setError = (errMsg) => {
    this.form.meta.error = errMsg
  }
}

export default FormStore