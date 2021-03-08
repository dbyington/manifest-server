import RequestForm from "./RequestForm";
import Show from "./Show";

import './Current.css';

function Current(props) {
  return (
    <div className="Current">
      <p>To generate manifest data enter the URL, beginning with <q>https://</q>, below then select the hash type and click Submit.</p>
      <RequestForm makeRequest={props.makeRequest}/>
      <Show manifest={props.manifest} loading={props.loading}/>
    </div>
  );
}

export default Current;

