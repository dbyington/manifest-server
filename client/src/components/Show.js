import './Show.css';

function Show(props) {
  let manifest;
  if (props.manifest && props.manifest.id) {
    manifest = (<div><h1>Built Manifest</h1>
      <pre>{JSON.stringify(JSON.parse(atob(props.manifest.asEncodedJson)), null, "    ")},</pre></div>)
  }

  if (props.loading) {
    return (
      <div className="Show"><h2>Reading package.... this could take a minute or so... maybe get some coffee, go for
        a run... :-)</h2></div>
    );
  }

  return (
    <div className="Show">
      {manifest}
    </div>
  );
}

export default Show;


