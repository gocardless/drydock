<img src="https://i.imgur.com/8bV9SEO.png" width="551" height="124" alt="DryDock">

DryDock
=======

DryDock is a utility, intended to be run as a cron job, to clean up old and unused docker images. This is useful on build servers and deployment nodes where image turnover can be high.

Usage
-----

    $ drydock --help
    DryDock 0.0.1
    usage: drydock [options]

    Options:
      --dry-run                          don't delete images
      --age      <48h>                   delete images older than age
      --pattern  <^.*$>                  pattern for images to be deleted
      --docker   <tcp://127.0.0.1:2375>  docker host endpoint
