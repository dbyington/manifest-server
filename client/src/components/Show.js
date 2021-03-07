import './Show.css';

function Show(props) {
  const manifest = props.manifest ? <div><h1>Built Manifest</h1>
    <pre>{JSON.stringify(JSON.parse(atob(props.manifest.asEncodedJson)), null, "    ")},</pre></div> : '';

  return (
    <div className="Show">
      {manifest}
    </div>
  );
}

export default Show;


