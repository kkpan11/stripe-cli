#!/bin/sh

PACKAGE_MANAGER=${1:-}

if [ $# -eq 0 ]; then
  echo "Error! Missing package manager argument"
  exit 1
fi

run_install() {
    case $PACKAGE_MANAGER in
    homebrew)
        brew install stripe/stripe-cli/stripe
    ;;

    apt)
        curl -s https://packages.stripe.dev/api/security/keypair/stripe-cli-gpg/public | gpg --dearmor | sudo tee /usr/share/keyrings/stripe.gpg
        echo "deb [signed-by=/usr/share/keyrings/stripe.gpg] https://packages.stripe.dev/stripe-cli-debian-local stable main" | sudo tee -a /etc/apt/sources.list.d/stripe.list
        sudo apt update
        sudo apt install stripe
    ;;

    yum)
        yum -y install stripe
    ;;

    scoop)
        scoop bucket add stripe https://github.com/stripe/scoop-stripe-cli.git
        scoop install stripe
    ;;

    docker)
    ;;

    *)
        echo "Error! Invalid package manager supplied"
        echo ""
        echo_help
        exit 1
        ;;
    esac

    stripe --version
}

trigger_pagerduty_alert() {
    sh -c "$(curl -sL https://raw.githubusercontent.com/martindstone/pagerduty-cli/master/install.sh)"
    pd event alert --routing_key "$PAGERDUTY_INTEGRATION_KEY" \
    --summary "Failed to install Stripe CLI on one or more operating systems. Investigate here: https://github.com/stripe/stripe-cli/actions/workflows/install-test.yml" \
    --timestamp "\"$(date)\"" \
    --source "Stripe CLI GitHub Actions" \
    --severity critical
}

if ! run_install
then
    echo "Install failed. Retrying in 30 seconds..."
    sleep 30
    if ! run_install
    then
        echo "Install failed again. Retrying in 60 seconds..."
        sleep 60
        if ! run_install
        then
            echo "Install failed again. Retrying in 120 seconds..."
            sleep 120
            if ! run_install
            then
                echo "Install failed again. Retrying for the last time in 180 seconds..."
                sleep 180
                run_install
                if ! run_install
                then
                trigger_pagerduty_alert
                exit 1
                fi
            fi
        fi
    fi
fi
