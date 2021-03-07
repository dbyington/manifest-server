import { Form, FormField } from 'react-hooks-form'

import './Form.css';


export default function RequestForm(props) {
  const { makeRequest } = props;

  async function handleSubmit(data) {
    console.log(data);
    if (!data.pkgUrl.startsWith("https")) {
      alert('Package URL must begin with "https://"')
      return;
    }
    if (!data.pkgUrl.endsWith(".pkg")) {
      alert('Expected package to end with ".pkg"')
      return;
    }
    if (!data.hashType) {
      alert('Please select a hash type')
      return
    }

    await makeRequest(data)
  }

  return (
    <Form className="Form" name="requestForm" onSubmit={handleSubmit}>
      <FormField component="input" type="text"  className="pkgUrl" name="pkgUrl" placeholder="https url to pkg file" />
      <FormField component="input" type="radio" name="hashType" value="md5" />md5
      <FormField component="input" type="radio" name="hashType" value="sha256" />sha256
      <input type="submit" name="submit" />
    </Form>
  );
}
